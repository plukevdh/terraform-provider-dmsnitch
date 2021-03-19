package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/repaygithub/terraform-provider-dmsnitch/dmsnitch"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dmsnitch.Provider})
}
