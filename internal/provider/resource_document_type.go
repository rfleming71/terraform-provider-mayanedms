package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceDocumentType() *schema.Resource {
	return &schema.Resource{
		Create: resourceDocumentTypeCreate,
		Read:   resourceDocumentTypeRead,
		Update: resourceDocumentTypeUpdate,
		Delete: resourceDocumentTypeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDocumentTypeImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_time_period": {
				Type:     schema.TypeInt,
				Default:  30,
				Optional: true,
			},
			"delete_time_unit": {
				Type:         schema.TypeString,
				Default:      "days",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"minutes", "hours", "days"}, false),
			},
			"trash_time_period": {
				Type:     schema.TypeInt,
				Default:  nil,
				Optional: true,
			},
			"trash_time_unit": {
				Type:         schema.TypeString,
				Default:      nil,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"minutes", "hours", "days"}, false),
			},
			"filename_generator_backend": {
				Type:         schema.TypeString,
				Default:      "uuid",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"original", "uuid", "uuid_plus_original"}, false),
			},
			"filename_generator_backend_arguments": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
		},
	}
}

func resourceDocumentTypeCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newDocType := dataToDocumentType(d)

	docType, err := c.CreateDocumentType(*newDocType)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", docType.ID))

	return resourceDocumentTypeRead(d, m)
}

func resourceDocumentTypeRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	docType, err := c.GetDocumentTypeById(id)
	if err != nil {
		return err
	}

	return documentTypeToData(docType, d)
}

func resourceDocumentTypeUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	docType := dataToDocumentType(d)
	docType.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateDocumentType(*docType)
	if err != nil {
		return err
	}

	return resourceDocumentTypeRead(d, m)
}

func resourceDocumentTypeDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteDocumentType(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceDocumentTypeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	docType, err := c.GetDocumentTypeById(id)
	if err != nil {
		return rd, err
	}

	err = documentTypeToData(docType, d)
	return rd, err
}

func documentTypeToData(docType *client.DocumentType, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", docType.ID))
	if err := d.Set("label", docType.Label); err != nil {
		return err
	}
	if err := d.Set("delete_time_period", docType.DeleteTimePeriod); err != nil {
		return err
	}
	if err := d.Set("delete_time_unit", docType.DeleteTimeUnit); err != nil {
		return err
	}
	if err := d.Set("filename_generator_backend", docType.FileNameGeneratorBackend); err != nil {
		return err
	}
	if err := d.Set("filename_generator_backend_arguments", docType.FileNameGeneratorBackendArguments); err != nil {
		return err
	}
	if docType.TrashTimePeriod != nil {
		if err := d.Set("trash_time_period", docType.TrashTimePeriod); err != nil {
			return err
		}
	}
	if docType.TrashTimeUnit != nil {
		if err := d.Set("trash_time_unit", docType.TrashTimeUnit); err != nil {
			return err
		}

	}

	return nil
}

func dataToDocumentType(d *schema.ResourceData) *client.DocumentType {
	newDocType := client.DocumentType{
		Label:                             d.Get("label").(string),
		DeleteTimePeriod:                  d.Get("delete_time_period").(int),
		DeleteTimeUnit:                    d.Get("delete_time_unit").(string),
		FileNameGeneratorBackend:          d.Get("filename_generator_backend").(string),
		FileNameGeneratorBackendArguments: d.Get("filename_generator_backend_arguments").(string),
	}

	x := d.Get("trash_time_period").(int)
	if x > 0 {
		newDocType.TrashTimePeriod = &x
	}

	y := d.Get("trash_time_unit").(string)
	if y != "" {
		newDocType.TrashTimeUnit = &y
	}

	return &newDocType
}
