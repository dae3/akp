package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func (a *App) GetPublicIp() (ip string, err error) {
	c := http.Client{}
	r, err := c.Get("https://api.ipify.org?format=text")
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unexpected HTTP status %s from IPIfy API", r.Status)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	ip = string(body)
	return
}

func (a *App) CreateNSG(ip string) (err error) {
	c, err := armnetwork.NewSecurityGroupsClient(a.Subscription, a.AzID, nil)
	if err != nil {
		a.setMessage(fmt.Sprintf("Error creating NSG client: %s", err.Error()))
		return
	}

	_, err = c.BeginCreateOrUpdate(a.ctx, "foo", "foo",
		armnetwork.SecurityGroup{
			Location: a.Location,
			Properties: &armnetwork.SecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.Ptr("rule1"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
							DestinationAddressPrefix: to.Ptr("*"),
							DestinationPortRange:     to.Ptr("22"),
							Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
							Priority:                 to.Ptr[int32](130),
							SourceAddressPrefix:      to.Ptr(ip),
							SourcePortRange:          to.Ptr("*"),
							Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						},
					}},
			},
		}, nil)
	if err != nil {
		a.setMessage(fmt.Sprintf("Error creating NSG: %s", err.Error()))
		return
	}

	a.setMessage("Created NSG")
	return
}
