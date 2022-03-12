package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTagCreate,
		Read:   resourceTagRead,
		Update: resourceTagUpdate,
		Delete: resourceTagDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTagImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Description: "Short text used as the tag name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"color": {
				Description: "The RGB color values for the tag.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceTagCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newTag := dataToTag(d)

	tag, err := c.CreateTag(*newTag)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", tag.ID))

	return resourceTagRead(d, m)
}

func resourceTagRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetTagById(id)
	if err != nil {
		return err
	}

	return tagToData(source, d)
}

func resourceTagUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	tag := dataToTag(d)
	tag.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateTag(*tag)
	if err != nil {
		return err
	}

	return resourceTagRead(d, m)
}

func resourceTagDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteTag(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceTagImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	tag, err := c.GetTagById(id)
	if err != nil {
		return rd, err
	}

	err = tagToData(tag, d)
	return rd, err
}

func tagToData(tag *client.Tag, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", tag.ID))
	if err := d.Set("label", tag.Label); err != nil {
		return err
	}
	if err := d.Set("color", tag.Color); err != nil {
		return err
	}

	return nil
}

func dataToTag(d *schema.ResourceData) *client.Tag {
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.Tag{
		ID:    id,
		Label: d.Get("label").(string),
		Color: d.Get("color").(string),
	}

	return &newDocType
}
