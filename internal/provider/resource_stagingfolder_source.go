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

type stagingFolderSourceBackendDataType struct {
	FolderPath        string `json:"folder_path"`
	PreviewWidth      int    `json:"preview_width"`
	PreviewHeight     int    `json:"preview_height"`
	DeleteAfterUpload bool   `json:"delete_after_upload"`
	Uncompress        string `json:"uncompress"`
}

func resourceStagingFolderSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceStagingFolderSourceCreate,
		Read:   resourceStagingFolderSourceRead,
		Update: resourceStagingFolderSourceUpdate,
		Delete: resourceStagingFolderSourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStagingFolderSourceImport,
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
			"folder_path": {
				Description: "Server side filesystem path.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"preview_width": {
				Description: "Width value to be passed to the converter backend.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"preview_height": {
				Description: "Height value to be passed to the converter backend.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"delete_after_upload": {
				Description: "Delete the file after is has been successfully uploaded.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"uncompress": {
				Description:  "Whether to expand or not compressed archives.",
				Type:         schema.TypeString,
				Default:      "ask",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ask", "yes", "no"}, false),
			},
		},
	}
}

func resourceStagingFolderSourceCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newSource := dataToStagingFolderSource(d)

	source, err := c.CreateSource(*newSource)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", source.ID))

	return resourceStagingFolderSourceRead(d, m)
}

func resourceStagingFolderSourceRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetSourceById(id)
	if err != nil {
		return err
	}

	return stagingFolderSourceToData(source, d)
}

func resourceStagingFolderSourceUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	source := dataToStagingFolderSource(d)
	source.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateSource(*source)
	if err != nil {
		return err
	}

	return resourceStagingFolderSourceRead(d, m)
}

func resourceStagingFolderSourceDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteSource(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceStagingFolderSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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

	if source.BackendPath != "mayan.apps.sources.source_backends.staging_folder_backends.SourceBackendStagingFolder" {
		return rd, errors.New("identified source is not of type watch folder")
	}

	err = stagingFolderSourceToData(source, d)
	return rd, err
}

func stagingFolderSourceToData(source *client.Source, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", source.ID))
	if err := d.Set("label", source.Label); err != nil {
		return err
	}
	if err := d.Set("enabled", source.Enabled); err != nil {
		return err
	}

	var backendData stagingFolderSourceBackendDataType
	_ = json.Unmarshal([]byte(source.BackendData), &backendData)

	uncompress := uncompressedMapping[backendData.Uncompress]

	if err := d.Set("uncompress", uncompress); err != nil {
		return err
	}

	if err := d.Set("folder_path", backendData.FolderPath); err != nil {
		return err
	}

	if err := d.Set("preview_width", backendData.PreviewWidth); err != nil {
		return err
	}

	if err := d.Set("preview_height", backendData.PreviewHeight); err != nil {
		return err
	}

	if err := d.Set("delete_after_upload", backendData.DeleteAfterUpload); err != nil {
		return err
	}

	return nil
}

func dataToStagingFolderSource(d *schema.ResourceData) *client.Source {
	backendData, _ := json.Marshal(stagingFolderSourceBackendDataType{
		Uncompress:        strings.ToLower(string(d.Get("uncompress").(string)[:1])),
		FolderPath:        d.Get("folder_path").(string),
		PreviewWidth:      d.Get("preview_width").(int),
		PreviewHeight:     d.Get("preview_height").(int),
		DeleteAfterUpload: d.Get("delete_after_upload").(bool),
	})
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.Source{
		ID:          id,
		Label:       d.Get("label").(string),
		Enabled:     d.Get("enabled").(bool),
		BackendPath: "mayan.apps.sources.source_backends.staging_folder_backends.SourceBackendStagingFolder",
		BackendData: string(backendData),
	}

	return &newDocType
}
