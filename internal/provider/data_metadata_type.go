package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func dataMetadataType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataMetadataTypeRead,

		Schema: map[string]*schema.Schema{
			"label": {
				Description: "Short description of this metadata type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name used by other apps to reference this metadata type. Do not use python reserved words, or spaces.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"default": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lookup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parser": {
				Description: "The parser will reformat the value entered to conform to the expected format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"validator": {
				Description: "The validator will reject data entry if the value entered does not conform to the expected format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataMetadataTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.MayanEdmsClient)
	name, _ := d.Get("name").(string)

	source, err := c.GetMetadataTypeByName(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if source == nil {
		return diag.Errorf("Unable to find metadata type with name '%s'", name)
	}

	return dataMetadataTypeToData(source, d)
}

func dataMetadataTypeToData(metadataType *client.MetadataType, d *schema.ResourceData) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%v", metadataType.ID))
	if err := d.Set("label", metadataType.Label); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", metadataType.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("default", metadataType.Default); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lookup", metadataType.Lookup); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parser", metadataType.Parser); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("validator", metadataType.Validator); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
