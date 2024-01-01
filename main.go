package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	_ "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "akp",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domready,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
	/*
		azfactory, err := armresources.NewClientFactory(subscription, id, nil)
		if err != nil {
			panic(err)
		}

		rgclient := azfactory.NewResourceGroupsClient()

		rgname := "foo"
		l := location
		_, err = rgclient.CreateOrUpdate(context.Background(), rgname, armresources.ResourceGroup{Location: &l}, &armresources.ResourceGroupsClientCreateOrUpdateOptions{})
		if err != nil {
			panic(err)
		}

			compfactory, err := armcompute.NewClientFactory(subscription, id, nil)
			if err != nil {
				panic(err)
			}

			vmclient := compfactory.NewVirtualMachinesClient()

			vm, err := vmclient.BeginCreateOrUpdate(context.Background(), rg.Name, "foo", armcompute.VirtualMachine{Location: &l, })
	*/
}