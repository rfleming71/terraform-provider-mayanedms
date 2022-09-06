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

type watchFolderSourceBackendDataType struct {
	FolderPath           string `json:"folder_path"`
	InclueSubdirectories bool   `json:"include_subdirectories"`
	DocumentTypeId       int    `json:"document_type_id"`
	Interval             int    `json:"interval"`
	Uncompress           string `json:"uncompress"`
}

func resourceWatchFolderSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceWatchFolderSourceCreate,
		Read:   resourceWatchFolderSourceRead,
		Update: resourceWatchFolderSourceUpdate,
		Delete: resourceWatchFolderSourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWatchFolderSourceImport,
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
				Description: "Server side filesystem path to scan for files.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"include_subdirectories": {
				Description: "If enabled, not only will the folder path be scanned for files but also its subdirectories.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"document_type_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Required: true,
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

func resourceWatchFolderSourceCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newSource := dataToWatchFolderSource(d)

	source, err := c.CreateSource(*newSource)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", source.ID))

	return resourceWatchFolderSourceRead(d, m)
}

func resourceWatchFolderSourceRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetSourceById(id)
	if err != nil {
		return err
	}

	return watchFolderSourceToData(source, d)
}

func resourceWatchFolderSourceUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	source := dataToWatchFolderSource(d)
	source.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateSource(*source)
	if err != nil {
		return err
	}

	return resourceWatchFolderSourceRead(d, m)
}

func resourceWatchFolderSourceDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteSource(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceWatchFolderSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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

	if source.BackendPath != "mayan.apps.sources.source_backends.watch_folder_backends.SourceBackendWatchFolder" {
		return rd, errors.New("identified source is not of type watch folder")
	}

	err = watchFolderSourceToData(source, d)
	return rd, err
}

func watchFolderSourceToData(source *client.Source, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", source.ID))
	if err := d.Set("label", source.Label); err != nil {
		return err
	}
	if err := d.Set("enabled", source.Enabled); err != nil {
		return err
	}

	var backendData watchFolderSourceBackendDataType
	_ = json.Unmarshal([]byte(source.BackendData), &backendData)

	uncompress := uncompressedMapping[backendData.Uncompress]

	if err := d.Set("uncompress", uncompress); err != nil {
		return err
	}

	if err := d.Set("folder_path", backendData.FolderPath); err != nil {
		return err
	}

	if err := d.Set("include_subdirectories", backendData.InclueSubdirectories); err != nil {
		return err
	}

	if err := d.Set("document_type_id", backendData.DocumentTypeId); err != nil {
		return err
	}

	if err := d.Set("interval", backendData.Interval); err != nil {
		return err
	}

	return nil
}

func dataToWatchFolderSource(d *schema.ResourceData) *client.Source {
	backendData, _ := json.Marshal(watchFolderSourceBackendDataType{
		Uncompress:           strings.ToLower(string(d.Get("uncompress").(string)[:1])),
		FolderPath:           d.Get("folder_path").(string),
		InclueSubdirectories: d.Get("include_subdirectories").(bool),
		DocumentTypeId:       d.Get("document_type_id").(int),
		Interval:             d.Get("interval").(int),
	})
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.Source{
		ID:          id,
		Label:       d.Get("label").(string),
		Enabled:     d.Get("enabled").(bool),
		BackendPath: "mayan.apps.sources.source_backends.watch_folder_backends.SourceBackendWatchFolder",
		BackendData: string(backendData),
	}

	return &newDocType
}
