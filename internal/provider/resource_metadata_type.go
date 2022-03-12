package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceMetadataType() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetadataTypeCreate,
		Read:   resourceMetadataTypeRead,
		Update: resourceMetadataTypeUpdate,
		Delete: resourceMetadataTypeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMetadataTypeImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Description: "Short description of this metadata type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name used by other apps to reference this metadata type. Do not use python reserved words, or spaces.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"default": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"lookup": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"parser": {
				Description: "The parser will reformat the value entered to conform to the expected format.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"validator": {
				Description: "The validator will reject data entry if the value entered does not conform to the expected format.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func resourceMetadataTypeCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newMetadataType := dataToMetadataType(d)

	metadataType, err := c.CreateMetadataType(*newMetadataType)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", metadataType.ID))

	return resourceMetadataTypeRead(d, m)
}

func resourceMetadataTypeRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetMetadataTypeById(id)
	if err != nil {
		return err
	}

	return metadataTypeToData(source, d)
}

func resourceMetadataTypeUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	metadataType := dataToMetadataType(d)
	metadataType.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateMetadataType(*metadataType)
	if err != nil {
		return err
	}

	return resourceMetadataTypeRead(d, m)
}

func resourceMetadataTypeDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteMetadataType(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceMetadataTypeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	metadataType, err := c.GetMetadataTypeById(id)
	if err != nil {
		return rd, err
	}

	err = metadataTypeToData(metadataType, d)
	return rd, err
}

func metadataTypeToData(metadataType *client.MetadataType, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", metadataType.ID))
	if err := d.Set("label", metadataType.Label); err != nil {
		return err
	}
	if err := d.Set("name", metadataType.Name); err != nil {
		return err
	}

	if err := d.Set("default", metadataType.Default); err != nil {
		return err
	}

	if err := d.Set("lookup", metadataType.Lookup); err != nil {
		return err
	}

	if err := d.Set("parser", metadataType.Parser); err != nil {
		return err
	}

	if err := d.Set("validator", metadataType.Validator); err != nil {
		return err
	}

	return nil
}

func dataToMetadataType(d *schema.ResourceData) *client.MetadataType {
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.MetadataType{
		ID:        id,
		Label:     d.Get("label").(string),
		Name:      d.Get("name").(string),
		Default:   d.Get("default").(string),
		Lookup:    d.Get("lookup").(string),
		Parser:    d.Get("parser").(string),
		Validator: d.Get("validator").(string),
	}

	return &newDocType
}
