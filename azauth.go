package main

import (
	"context"

	_ "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

const (
	subscription = "86752990-af4d-4016-abac-65a9c2f536ac"
	location     = "australiaeast"
)

func (app *App) AzAuthenticate() (err error) {
	app.AzID, err = azidentity.NewDefaultAzureCredential(nil)
	return
}

type Subscription struct {
	DisplayName string
	Id          string
}

func (app *App) GetSubscriptions() (subscriptions []Subscription, err error) {
	tc, err := armsubscriptions.NewClient(app.AzID, nil)
	if err != nil {
		return
	}

	subspager := tc.NewListPager(nil)
	for subspager.More() {
		var page armsubscriptions.ClientListResponse
		page, err = subspager.NextPage(context.Background())
		for _, v := range page.Value {
			subscriptions = append(subscriptions, Subscription{DisplayName: *v.DisplayName, Id: *v.SubscriptionID})
		}
	}

	// default to first one
	app.Subscription = subscriptions[0].Id
	return
}

func (app *App) SetSubscription(sub string) {
	app.Subscription = sub
}
