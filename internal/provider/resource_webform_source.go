package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

type webformSourceBackendDataType struct {
	Uncompress string `json:"uncompress"`
}

var uncompressedMapping = map[string]string{
	"y": "yes",
	"n": "no",
	"a": "ask",
}

func resourceWebformSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebformSourceCreate,
		Read:   resourceWebformSourceRead,
		Update: resourceWebformSourceUpdate,
		Delete: resourceWebformSourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWebformSourceImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"uncompress": {
				Type:         schema.TypeString,
				Default:      "ask",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ask", "yes", "no"}, false),
			},
		},
	}
}

func resourceWebformSourceCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newSource := dataToWebformSource(d)

	source, err := c.CreateSource(*newSource)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", source.ID))

	return resourceWebformSourceRead(d, m)
}

func resourceWebformSourceRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetSourceById(id)
	if err != nil {
		return err
	}

	return webformSourceToData(source, d)
}

func resourceWebformSourceUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	source := dataToWebformSource(d)
	source.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateSource(*source)
	if err != nil {
		return err
	}

	return resourceWebformSourceRead(d, m)
}

func resourceWebformSourceDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteSource(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceWebformSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	source, err := c.GetSourceById(id)
	if err != nil {
		return rd, err
	}

	if source.BackendPath != "mayan.apps.sources.source_backends.SourceBackendWebForm" {
		return rd, errors.New("identified source is not of type webform")
	}

	err = webformSourceToData(source, d)
	return rd, err
}

func webformSourceToData(source *client.Source, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", source.ID))
	if err := d.Set("label", source.Label); err != nil {
		return err
	}
	if err := d.Set("enabled", source.Enabled); err != nil {
		return err
	}

	var backendData webformSourceBackendDataType
	_ = json.Unmarshal([]byte(source.BackendData), &backendData)

	uncompress := uncompressedMapping[backendData.Uncompress]

	if err := d.Set("uncompress", uncompress); err != nil {
		return err
	}

	return nil
}

func dataToWebformSource(d *schema.ResourceData) *client.Source {
	backendData, _ := json.Marshal(webformSourceBackendDataType{
		Uncompress: strings.ToLower(string(d.Get("uncompress").(string)[:1])),
	})
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.Source{
		ID:          id,
		Label:       d.Get("label").(string),
		Enabled:     d.Get("enabled").(bool),
		BackendPath: "mayan.apps.sources.source_backends.web_form_backends.SourceBackendWebForm",
		BackendData: string(backendData),
	}

	return &newDocType
}
