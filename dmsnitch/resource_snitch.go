package dmsnitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

type Snitch struct {
	Token    string   `json:"token,omitempty"`
	URL      string   `json:"check_in_url,omitempty"`
	Name     string   `json:"name,omitempty"`
	Status   string   `json:"status,omitempty"`
	Interval string   `json:"interval,omitempty"`
	Type     string   `json:"alert_type,omitempty"`
	Notes    string   `json:"notes,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func resourceSnitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnitchCreate,
		Update: resourceSnitchUpdate,
		Read:   resourceSnitchRead,
		Delete: resourceSnitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "basic",
			},

			"interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "daily",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func newSnitchFromResource(d *schema.ResourceData) *Snitch {
	tags := make([]string, 0, len(d.Get("tags").(*schema.Set).List()))

	for _, item := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, item.(string))
	}

	return &Snitch{
		Name:     d.Get("name").(string),
		Token:    d.Get("token").(string),
		URL:      d.Get("url").(string),
		Status:   d.Get("status").(string),
		Interval: d.Get("interval").(string),
		Type:     d.Get("type").(string),
		Notes:    d.Get("notes").(string),
		Tags:     tags,
	}
}

func resourceSnitchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	snitch := newSnitchFromResource(d)

	bytedata, err := json.Marshal(snitch)
	if err != nil {
		return err
	}

	resp, err := client.Post("snitches", bytes.NewBuffer(bytedata)) //nolint:bodyclose
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		body, readerr := ioutil.ReadAll(resp.Body)

		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &snitch)

		if decodeerr != nil {
			return decodeerr
		}

		log.Printf("[DEBUG] ID received: %s", snitch.Token)
		d.SetId(snitch.Token)
	}

	return resourceSnitchRead(d, m)
}

func resourceSnitchRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	resp, _ := client.Get(fmt.Sprintf("snitches/%s", d.Id())) //nolint:bodyclose

	if resp.StatusCode == http.StatusOK { //nolint:nestif
		var snitch Snitch

		body, readerr := ioutil.ReadAll(resp.Body)

		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &snitch)

		if decodeerr != nil {
			return decodeerr
		}

		tagList := make([]string, 0, len(snitch.Tags))
		tagList = append(tagList, snitch.Tags...)
		if err := d.Set("name", snitch.Name); err != nil {
			return err
		}
		if err := d.Set("token", snitch.Token); err != nil {
			return err
		}
		if err := d.Set("url", snitch.URL); err != nil {
			return err
		}
		if err := d.Set("status", snitch.Status); err != nil {
			return err
		}
		if err := d.Set("interval", snitch.Interval); err != nil {
			return err
		}
		if err := d.Set("type", snitch.Type); err != nil {
			return err
		}
		if err := d.Set("notes", snitch.Notes); err != nil {
			return err
		}
		if err := d.Set("tags", tagList); err != nil {
			return err
		}
	}

	return nil
}

func resourceSnitchUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	snitch := newSnitchFromResource(d)

	var jsonBuffer []byte

	jsonPayload := bytes.NewBuffer(jsonBuffer)
	enc := json.NewEncoder(jsonPayload)
	if err := enc.Encode(snitch); err != nil {
		return err
	}

	id := d.Id()

	_, err := client.Patch(fmt.Sprintf("snitches/%s", id), jsonPayload) //nolint:bodyclose
	if err != nil {
		return err
	}

	return resourceSnitchRead(d, m)
}

func resourceSnitchDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf("snitches/%s", id)) //nolint:bodyclose

	return err
}
