package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Resource struct {
	Name   string
	Status string
}

// App struct
type App struct {
	ctx          context.Context
	AzID         *azidentity.DefaultAzureCredential
	Resources    []*Resource
	Location     string
	Subscription string
}

// NewApp creates a new App application struct
func NewApp() (a *App) {
	a = &App{}
	a.Location = "australiaeast"
	a.Resources = []*Resource{
		&Resource{"Resource Group", "Unknown"},
		&Resource{"NSG", "Unknown"},
		&Resource{"Proxy VM", "Unknown"},
		&Resource{"kubeconfig", "Unknown"},
	}
	return
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = a.AzAuthenticate()
}

func (a *App) domready(ctx context.Context) {

	s, e := a.GetSubscriptions()
	if e != nil {
		print(e)
	}
	runtime.EventsEmit(ctx, "resource_change", a.Resources)
	runtime.EventsEmit(ctx, "subscriptions_loaded", s)
	runtime.EventsEmit(ctx, "message", "hi")
}

func (a *App) Connect() {
	a.setMessage("Connecting...")
	if err := a.CheckOrCreateResourceGroup("foo"); err != nil {
		a.setMessage(fmt.Sprintf("Resource Group creation failed: %s", err.Error()))
		return
	}
	a.setMessage("Resource Group created")
	for _, r := range a.Resources {
		if r.Name == "Resource Group" {
			r.Status = "OK"
		}
	}
	runtime.EventsEmit(a.ctx, "resource_change", a.Resources)

	ip, err := a.GetPublicIp()
	if err != nil {
		a.setMessage(fmt.Sprintf("Error determining public IP: %s", err.Error()))
		return
	}
	nsg, err := a.CreateNSG(ip)
	if err != nil {
		a.setMessage(fmt.Sprintf("NSG creation failed: %s", err.Error()))
		return
	}
	for _, r := range a.Resources {
		if r.Name == "NSG" {
			r.Status = fmt.Sprintf("OK (your public IP is %s)", ip)
		}
	}
	runtime.EventsEmit(a.ctx, "resource_change", a.Resources)

	a.setMessage("Creating Proxy VM...")
	if err = a.CreateProxyVM("foo", nsg); err != nil {
		a.setMessage(fmt.Sprintf("Proxy VM creation failed: %s", err.Error()))
		return
	}
	for _, r := range a.Resources {
		if r.Name == "Proxy VM" {
			r.Status = fmt.Sprintf("OK", ip)
		}
	}

	a.setMessage("")
}

func (a *App) setMessage(msg string) {
	runtime.EventsEmit(a.ctx, "message", msg)
}
