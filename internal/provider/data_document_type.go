package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func dataDocumentType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataDocumentTypeRead,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_time_period": {
				Description: "Amount of time after which documents of this type will be moved to the trash",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"delete_time_unit": {
				Description:  "Unit of delete_time_period. (minutes, hours, days)",
				Type:         schema.TypeString,
				Default:      "days",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"minutes", "hours", "days"}, false),
			},
			"trash_time_period": {
				Description: "Amount of time after which documents of this type in the trash will be deleted.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"trash_time_unit": {
				Description: "Unit of trash_time_period. (minutes, hours, days)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"filename_generator_backend": {
				Description: "The class responsible for producing the actual filename used to store the uploaded documents",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"filename_generator_backend_arguments": {
				Description: "The arguments for the filename generator backend as a YAML dictionary.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataDocumentTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.MayanEdmsClient)
	label, _ := d.Get("label").(string)

	docTypes, err := c.GetDocumentTypes()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, docType := range docTypes {
		if docType.Label == label {
			return dataDocumentTypeToData(&docType, d)
		}
	}

	return diag.Errorf("Unable to find document type with label '%s'", label)
}

func dataDocumentTypeToData(docType *client.DocumentType, d *schema.ResourceData) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%v", docType.ID))
	if err := d.Set("label", docType.Label); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("delete_time_period", docType.DeleteTimePeriod); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("delete_time_unit", docType.DeleteTimeUnit); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("filename_generator_backend", docType.FileNameGeneratorBackend); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("filename_generator_backend_arguments", docType.FileNameGeneratorBackendArguments); err != nil {
		return diag.FromErr(err)
	}
	if docType.TrashTimePeriod != nil {
		if err := d.Set("trash_time_period", docType.TrashTimePeriod); err != nil {
			return diag.FromErr(err)
		}
	}
	if docType.TrashTimeUnit != nil {
		if err := d.Set("trash_time_unit", docType.TrashTimeUnit); err != nil {
			return diag.FromErr(err)
		}

	}

	return nil
}
