package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceIndexTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceIndexTemplateCreate,
		Read:   resourceIndexTemplateRead,
		Update: resourceIndexTemplateUpdate,
		Delete: resourceIndexTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIndexTemplateImport,
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
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"document_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceIndexTemplateCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newIndexTemplate := dataToIndexTemplate(d)

	indexTemplate, err := c.CreateIndexTemplate(*newIndexTemplate)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", indexTemplate.ID))

	dTypes := d.Get("document_types").(*schema.Set).List()
	documentTypes := []int{}
	for _, tag := range dTypes {
		t := tag.(int)
		documentTypes = append(documentTypes, t)
	}

	for _, docType := range documentTypes {
		c.AddIndexTemplateDocumentType(indexTemplate.ID, docType)
	}

	return resourceIndexTemplateRead(d, m)
}

func resourceIndexTemplateRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetIndexTemplateById(id)
	if err != nil {
		return err
	}

	err = indexTemplateToData(source, d)
	if err != nil {
		return err
	}

	docTypes, err := c.GetIndexTemplateDocumentTypes(source.ID)
	if err != nil {
		return err
	}

	if err := d.Set("document_types", docTypes); err != nil {
		return err
	}

	return nil
}

func resourceIndexTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	indexTemplate := dataToIndexTemplate(d)
	indexTemplate.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateIndexTemplate(*indexTemplate)
	if err != nil {
		return err
	}

	if d.HasChange("document_types") {
		o, n := d.GetChange("document_types")

		oTypes := o.(*schema.Set)
		nTypes := n.(*schema.Set)

		removals := oTypes.Difference(nTypes)
		for _, removal := range removals.List() {
			if err := c.RemoveIndexTemplateDocumentType(indexTemplate.ID, removal.(int)); err != nil {
				return err
			}
		}

		additions := nTypes.Difference(oTypes)
		for _, addition := range additions.List() {
			if err := c.AddIndexTemplateDocumentType(indexTemplate.ID, addition.(int)); err != nil {
				return err
			}
		}
	}

	return resourceIndexTemplateRead(d, m)
}

func resourceIndexTemplateDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteIndexTemplate(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceIndexTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	indexTemplate, err := c.GetIndexTemplateById(id)
	if err != nil {
		return rd, err
	}

	err = indexTemplateToData(indexTemplate, d)
	if err != nil {
		return rd, err
	}

	docTypes, err := c.GetIndexTemplateDocumentTypes(indexTemplate.ID)
	if err != nil {
		return rd, err
	}

	if err := d.Set("document_types", docTypes); err != nil {
		return rd, err
	}

	return rd, err
}

func indexTemplateToData(indexTemplate *client.IndexTemplate, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", indexTemplate.ID))
	if err := d.Set("label", indexTemplate.Label); err != nil {
		return err
	}
	if err := d.Set("slug", indexTemplate.Slug); err != nil {
		return err
	}
	if err := d.Set("enabled", indexTemplate.Enabled); err != nil {
		return err
	}

	return nil
}

func dataToIndexTemplate(d *schema.ResourceData) *client.IndexTemplate {
	id, _ := strconv.Atoi(d.Id())
	newIndexTemplate := client.IndexTemplate{
		ID:      id,
		Label:   d.Get("label").(string),
		Slug:    d.Get("slug").(string),
		Enabled: d.Get("enabled").(bool),
	}

	return &newIndexTemplate
}
