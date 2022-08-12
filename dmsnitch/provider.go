package dmsnitch

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
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
	client := &DMSnitchClient{
		ApiKey:     d.Get("api_key").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
