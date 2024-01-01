package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func (a *App) CheckOrCreateResourceGroup(name string) (err error) {
	factory, err := armresources.NewClientFactory(subscription, a.AzID, nil)
	if err != nil {
		return
	}
	rgclient := factory.NewResourceGroupsClient()

	_, err = rgclient.CreateOrUpdate(a.ctx, name, armresources.ResourceGroup{Location: to.Ptr(a.Location)}, nil)
	if err != nil {
		return
	}
	return nil
}
