package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceIndexTemplateNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceIndexTemplateNodeCreate,
		Read:   resourceIndexTemplateNodeRead,
		Update: resourceIndexTemplateNodeUpdate,
		Delete: resourceIndexTemplateNodeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIndexTemplateNodeImport,
		},

		Schema: map[string]*schema.Schema{
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				// Avoid issues due to trailing whitespace
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					suppressDiff := strings.TrimSpace(old) == strings.TrimSpace(new)
					return suppressDiff
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"link_documents": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"index_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"parent_id": {
				Type: schema.TypeInt,
				// TODO: until we are able to compute this
				Required: true,
				// Optional: true,
			},
		},
	}
}

func resourceIndexTemplateNodeCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newIndexTemplateNode := dataToIndexTemplateNode(d)

	indexTemplateNode, err := c.CreateIndexTemplateNode(*newIndexTemplateNode)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", indexTemplateNode.ID))

	return resourceIndexTemplateNodeRead(d, m)
}

func resourceIndexTemplateNodeRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	indexId := d.Get("index_id").(int)

	source, err := c.GetIndexTemplateNodeById(indexId, id)
	if err != nil {
		return err
	}

	err = indexTemplateNodeToData(source, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceIndexTemplateNodeUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	indexTemplateNode := dataToIndexTemplateNode(d)
	indexTemplateNode.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateIndexTemplateNode(*indexTemplateNode)
	if err != nil {
		return err
	}

	return resourceIndexTemplateNodeRead(d, m)
}

func resourceIndexTemplateNodeDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	indexId := d.Get("index_id").(int)
	err := c.DeleteIndexTemplateNode(indexId, id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceIndexTemplateNodeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	indexId := d.Get("index_id").(int)

	indexTemplateNode, err := c.GetIndexTemplateNodeById(indexId, id)
	if err != nil {
		return rd, err
	}

	err = indexTemplateNodeToData(indexTemplateNode, d)
	if err != nil {
		return rd, err
	}

	return rd, err
}

func indexTemplateNodeToData(indexTemplateNode *client.IndexTemplateNode, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", indexTemplateNode.ID))
	if err := d.Set("expression", indexTemplateNode.Expression); err != nil {
		return err
	}
	if err := d.Set("enabled", indexTemplateNode.Enabled); err != nil {
		return err
	}
	if err := d.Set("link_documents", indexTemplateNode.LinkDocuments); err != nil {
		return err
	}
	if err := d.Set("index_id", indexTemplateNode.IndexID); err != nil {
		return err
	}
	if err := d.Set("parent_id", indexTemplateNode.ParentID); err != nil {
		return err
	}

	return nil
}

func dataToIndexTemplateNode(d *schema.ResourceData) *client.IndexTemplateNode {
	id, _ := strconv.Atoi(d.Id())
	newIndexTemplateNode := client.IndexTemplateNode{
		ID:            id,
		Expression:    d.Get("expression").(string),
		Enabled:       d.Get("enabled").(bool),
		LinkDocuments: d.Get("link_documents").(bool),
		IndexID:       d.Get("index_id").(int),
		Parent:        d.Get("parent_id").(int),
		ParentID:      d.Get("parent_id").(int),
	}

	return &newIndexTemplateNode
}
