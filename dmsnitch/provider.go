package dmsnitch

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
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
	client := &DMSnitchClient{
		ApiKey:     d.Get("api_key").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
