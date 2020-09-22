package dmsnitch

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DMS_TOKEN", nil),
				Description: "Dead Man's Snitch API Key",
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"dmsnitch_snitch": resourceSnitch(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &Client{
		APIKey:     d.Get("api_key").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
