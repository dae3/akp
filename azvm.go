package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func (app *App) CreateProxyVM(ResourceGroupName string, nsg armnetwork.SecurityGroup) (err error) {
	vmclient, err := armcompute.NewVirtualMachinesClient(app.Subscription, app.AzID, nil)
	if err != nil {
		return
	}

	netclient, err := armnetwork.NewSubnetsClient(app.Subscription, app.AzID, nil)
	if err != nil {
		return
	}

	subnet, err := netclient.Get(app.ctx, "net", "net", "subnet", nil)
	if err != nil {
		return
	}

	sshkeys, err := GenerateSSHKeyPair()
	if err != nil {
		return
	}

	poller, err := vmclient.BeginCreateOrUpdate(
		app.ctx, ResourceGroupName, "foo",
		armcompute.VirtualMachine{
			Location: to.Ptr(app.Location),
			Properties: &armcompute.VirtualMachineProperties{
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: to.Ptr(armcompute.VirtualMachineSizeTypesStandardB1Ms),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkAPIVersion: to.Ptr(armcompute.NetworkAPIVersionTwoThousandTwenty1101),
					NetworkInterfaceConfigurations: []*armcompute.VirtualMachineNetworkInterfaceConfiguration{
						&armcompute.VirtualMachineNetworkInterfaceConfiguration{
							Name: to.Ptr("main"),
							Properties: &armcompute.VirtualMachineNetworkInterfaceConfigurationProperties{
								IPConfigurations: []*armcompute.VirtualMachineNetworkInterfaceIPConfiguration{
									&armcompute.VirtualMachineNetworkInterfaceIPConfiguration{
										Name: to.Ptr("main"),
										Properties: &armcompute.VirtualMachineNetworkInterfaceIPConfigurationProperties{
											PublicIPAddressConfiguration: &armcompute.VirtualMachinePublicIPAddressConfiguration{
												Name: to.Ptr("main"),
												Properties: &armcompute.VirtualMachinePublicIPAddressConfigurationProperties{
													DNSSettings: &armcompute.VirtualMachinePublicIPAddressDNSSettingsConfiguration{
														DomainNameLabel: to.Ptr("foo"),
													},
													DeleteOption:             to.Ptr(armcompute.DeleteOptionsDelete),
													PublicIPAllocationMethod: to.Ptr(armcompute.PublicIPAllocationMethodStatic),
												},
												SKU: &armcompute.PublicIPAddressSKU{Name: to.Ptr(armcompute.PublicIPAddressSKUNameStandard)},
											},
											Subnet: &armcompute.SubResource{ID: subnet.Subnet.ID},
										},
									},
								},
								NetworkSecurityGroup: &armcompute.SubResource{ID: nsg.ID},
							},
						},
					},
				},
				OSProfile: &armcompute.OSProfile{
					AdminUsername: to.Ptr("proxyadmin"),
					ComputerName:  to.Ptr("foo"),
					LinuxConfiguration: &armcompute.LinuxConfiguration{
						DisablePasswordAuthentication: to.Ptr(true),
						SSH: &armcompute.SSHConfiguration{
							PublicKeys: []*armcompute.SSHPublicKey{&armcompute.SSHPublicKey{
								KeyData: to.Ptr(sshkeys.Public),
								Path:    to.Ptr(fmt.Sprintf("/home/proxyadmin/.ssh/authorized_keys")),
							}},
						},
					},
				},
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Offer:     to.Ptr("0001-com-ubuntu-server-jammy"),
						Publisher: to.Ptr("Canonical"),
						Version:   to.Ptr("latest"),
						SKU:       to.Ptr("22_04-lts-gen2"),
					},
					OSDisk: &armcompute.OSDisk{
						CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
						DeleteOption: to.Ptr(armcompute.DiskDeleteOptionTypesDelete),
						DiskSizeGB:   to.Ptr(int32(30)),
					},
				},
			},
		},
		nil,
	)

	for {
		if poller.Done() {
			break
		}
		var resp *http.Response
		resp, err = poller.Poll(app.ctx)
		if err != nil {
			return
		}
		var msg []byte
		msg, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		app.setMessage(string(msg))
	}
	return
}
