package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			StateContext: resourceIndexTemplateNodeImport,
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
				Required: true,
			},
			"node_id": {
				Type: schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIndexTemplateNodeCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)

	indexTemplateId, newIndexTemplateNode := dataToIndexTemplateNode(d)

	indexTemplateNode, err := c.CreateIndexTemplateNode(*newIndexTemplateNode)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v-%v", indexTemplateId, indexTemplateNode.ID))

	return resourceIndexTemplateNodeRead(d, m)
}

func resourceIndexTemplateNodeRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	indexTemplateId, id, err := getIdInformation(d)
	if err != nil {
		return err
	}

	source, err := c.GetIndexTemplateNodeById(indexTemplateId, id)
	if err != nil {
		return err
	}

	err = indexTemplateNodeToData(indexTemplateId, source, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceIndexTemplateNodeUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	_, indexTemplateNode := dataToIndexTemplateNode(d)
	indexTemplateId, indexTemplateNodeId, err := getIdInformation(d)
	if err != nil {
		return err
	}

	indexTemplateNode.ID = indexTemplateNodeId
	_, err = c.UpdateIndexTemplateNode(indexTemplateId, *indexTemplateNode)
	if err != nil {
		return err
	}

	return resourceIndexTemplateNodeRead(d, m)
}

func resourceIndexTemplateNodeDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	indexTemplateId, id, err := getIdInformation(d)
	if err != nil {
		return err
	}

	err = c.DeleteIndexTemplateNode(indexTemplateId, id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceIndexTemplateNodeImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	indexTemplateId, id, err := getIdInformation(d)
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	tflog.Debug(ctx, "Requesting IndexTemplateNode by Id", map[string]interface{}{
		"index_id": indexTemplateId,
		"id": id,
	})
	indexTemplateNode, err := c.GetIndexTemplateNodeById(indexTemplateId, id)
	if err != nil {
		return rd, err
	}

	err = indexTemplateNodeToData(indexTemplateId, indexTemplateNode, d)
	if err != nil {
		return rd, err
	}

	return rd, err
}

func indexTemplateNodeToData(indexTemplateId int, indexTemplateNode *client.IndexTemplateNode, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v-%v", indexTemplateId, indexTemplateNode.ID))
	if err := d.Set("node_id", indexTemplateNode.ID); err != nil {
		return err
	}
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

func dataToIndexTemplateNode(d *schema.ResourceData) (int, *client.IndexTemplateNode) {
	_, id, _ := getIdInformation(d)
	newIndexTemplateNode := client.IndexTemplateNode{
		ID:            id,
		Expression:    d.Get("expression").(string),
		Enabled:       d.Get("enabled").(bool),
		LinkDocuments: d.Get("link_documents").(bool),
		IndexID:       d.Get("index_id").(int),
		Parent:        d.Get("parent_id").(int),
		ParentID:      d.Get("parent_id").(int),
	}

	return d.Get("index_id").(int), &newIndexTemplateNode
}
