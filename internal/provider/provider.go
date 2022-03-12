// provider.go
package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"url": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("MAYAN_EDMS_URL", nil),
					Description: "Hostname of the mayan edms host",
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("MAYAN_EDMS_USER", nil),
					Description: "User account for mayan edms api",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("MAYAN_EDMS_PASSWORD", nil),
					Description: "Password for mayan edms api",
				},
				"insecure": &schema.Schema{
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					DefaultFunc: schema.EnvDefaultFunc("MAYAN_EDMS_INSECURE", nil),
					Description: "Whether SSL should be verified or not",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"mayanedms_document_type":                resourceDocumentType(),
				"mayanedms_webform_source":               resourceWebformSource(),
				"mayanedms_watchfolder_source":           resourceWatchFolderSource(),
				"mayanedms_stagingfolder_source":         resourceStagingFolderSource(),
				"mayanedms_tag":                          resourceTag(),
				"mayanedms_index_template":               resourceIndexTemplate(),
				"mayanedms_group":                        resourceGroup(),
				"mayanedms_workflow_template":            resourceWorkflowTemplate(),
				"mayanedms_workflow_template_state":      resourceWorkflowTemplateState(),
				"mayanedms_workflow_template_transition": resourceWorkflowTemplateTransition(),
				"mayanedms_role":                         resourceRole(),
				"mayanedms_metadata_type":                resourceMetadataType(),
			},
			ConfigureFunc: mayanEdmsConfigure,
		}
		return p
	}
}

func mayanEdmsConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(bool)
	config := client.ClientConfig{
		Url:                url,
		Username:           username,
		Password:           password,
		InsecureSkipVerify: insecure,
	}
	return client.NewMayanEdmsClient(config)
}
