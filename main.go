package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/plukevdh/terraform-provider-dmsnitch/dmsnitch"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dmsnitch.Provider})
}
