package dmsnitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

type Snitch struct {
	Token    string   `json:"token,omitempty"`
	Url      string   `json:"check_in_url,omitempty"`
	Name     string   `json:"name,omitempty"`
	Status   string   `json:"status,omitempty"`
	Interval string   `json:"interval,omitempty"`
	Type     string   `json:"alert_type,omitempty"`
	Email    []string `json:"alert_email,omitempty"`
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

			"alert_email": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
	email := make([]string, 0, len(d.Get("alert_email").(*schema.Set).List()))

	for _, item := range d.Get("alert_email").(*schema.Set).List() {
		email = append(email, item.(string))
	}

	tags := make([]string, 0, len(d.Get("tags").(*schema.Set).List()))

	for _, item := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, item.(string))
	}

	return &Snitch{
		Name:     d.Get("name").(string),
		Token:    d.Get("token").(string),
		Url:      d.Get("url").(string),
		Status:   d.Get("status").(string),
		Interval: d.Get("interval").(string),
		Type:     d.Get("type").(string),
		Email:    email,
		Notes:    d.Get("notes").(string),
		Tags:     tags,
	}
}

func resourceSnitchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*DMSnitchClient)
	snitch := newSnitchFromResource(d)

	bytedata, err := json.Marshal(snitch)

	if err != nil {
		return err
	}

	resp, err := client.Post("snitches", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
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
	client := m.(*DMSnitchClient)
	resp, _ := client.Get(fmt.Sprintf("snitches/%s", d.Id()))

	if resp.StatusCode == 200 {
		var snitch Snitch

		body, readerr := ioutil.ReadAll(resp.Body)

		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &snitch)

		if decodeerr != nil {
			return decodeerr
		}

		emailList := make([]string, 0, len(snitch.Email))

		for _, event := range snitch.Email {
			emailList = append(emailList, event)
		}

		tagList := make([]string, 0, len(snitch.Tags))

		for _, event := range snitch.Tags {
			tagList = append(tagList, event)
		}

		d.Set("name", snitch.Name)
		d.Set("token", snitch.Token)
		d.Set("url", snitch.Url)
		d.Set("status", snitch.Status)
		d.Set("interval", snitch.Interval)
		d.Set("type", snitch.Type)
		d.Set("tags", emailList)
		d.Set("notes", snitch.Notes)
		d.Set("tags", tagList)
	}

	return nil
}

func resourceSnitchUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*DMSnitchClient)
	snitch := newSnitchFromResource(d)

	var jsonBuffer []byte

	jsonPayload := bytes.NewBuffer(jsonBuffer)
	enc := json.NewEncoder(jsonPayload)
	enc.Encode(snitch)

	id := d.Id()

	_, err := client.Patch(fmt.Sprintf("snitches/%s", id), jsonPayload)

	if err != nil {
		return err
	}

	return resourceSnitchRead(d, m)
}

func resourceSnitchDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*DMSnitchClient)
	_, err := client.Delete(fmt.Sprintf("snitches/%s", id))

	return err
}
