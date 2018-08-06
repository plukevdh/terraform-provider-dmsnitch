package dmsnitch

import (
	"github.com/hashicorp/terraform/helper/schema"
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
		"net/http"
)

type Snitch struct {
	Token    string   `json:"token,omitempty"`
	Href     string   `json:"href,omitempty"`
	Url      string   `json:"check_in_url,omitempty"`
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

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"alert_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "basic",
			},

			"interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "daily",
			},

			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
			},

			"url": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},

			"href": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func newSnitchFromResource(d *schema.ResourceData) *Snitch {
	snitch := &Snitch{
		Name:     d.Get("name").(string),
		Href:     d.Get("href").(string),
		Token:    d.Get("token").(string),
		Url:      d.Get("check_in_url").(string),
		Status:   d.Get("status").(string),
		Interval: d.Get("interval").(string),
		Type:     d.Get("type").(string),
		Notes:    d.Get("notes").(string),
		Tags:     d.Get("tags").([]string),
	}

	return snitch
}

func newResourceFromApi(resp *http.Response, d *schema.ResourceData) error {
	var snitch Snitch

	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &snitch)
	if decodeerr != nil {
		return decodeerr
	}
	
	d.SetId(snitch.Token)
	d.Set("name", snitch.Name)
	d.Set("href", snitch.Href)
	d.Set("token", snitch.Token)
	d.Set("url", snitch.Url)
	d.Set("status", snitch.Status)
	d.Set("interval", snitch.Interval)
	d.Set("type", snitch.Type)
	d.Set("notes", snitch.Notes)
	d.Set("tags", snitch.Tags)
	
	return nil;
}

func resourceSnitchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*DMSnitchClient)
	snitch := newSnitchFromResource(d)

	bytedata, err := json.Marshal(snitch)

	if err != nil {
		return err
	}

	body, err := client.Post("v1/snitches", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	if body.StatusCode == 201 {
		buildErr := newResourceFromApi(body, d)

		if buildErr != nil {
			return buildErr
		}
	}

	return nil
}

func resourceSnitchRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*DMSnitchClient)
	resp, _ := client.Get(fmt.Sprintf("v1/snitches/%s", id))

	if resp.StatusCode == 200 {
		buildErr := newResourceFromApi(resp, d)

		if buildErr != nil {
			return buildErr
		}
	}

	return nil
}

func resourceSnitchUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*DMSnitchClient)
	repository := newSnitchFromResource(d)

	var jsonBuffer []byte

	jsonPayload := bytes.NewBuffer(jsonBuffer)
	enc := json.NewEncoder(jsonPayload)
	enc.Encode(repository)

	id := d.Id()

	resp, _ := client.Patch(fmt.Sprintf("v1/snitches/%s", id), jsonPayload)

	if resp.StatusCode == 200 {
		buildErr := newResourceFromApi(resp, d)

		if buildErr != nil {
			return buildErr
		}
	}

	return nil
}

func resourceSnitchDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*DMSnitchClient)
	_, err := client.Delete(fmt.Sprintf("v1/snitches/%s", id))

	return err
}