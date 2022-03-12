package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceWorkflowTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkflowTemplateCreate,
		Read:   resourceWorkflowTemplateRead,
		Update: resourceWorkflowTemplateUpdate,
		Delete: resourceWorkflowTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWorkflowTemplateImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Description: "Short text to describe the workflow",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_name": {
				Description: "This value will be used by other apps to reference this workflow. Can only contain letters, numbers, and underscores.",
				Type:        schema.TypeString,
				Required:    true,
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

func resourceWorkflowTemplateCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newWorkflowTemplate := dataToWorkflowTemplate(d)

	workflowTemplate, err := c.CreateWorkflowTemplate(*newWorkflowTemplate)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", workflowTemplate.ID))

	dTypes := d.Get("document_types").(*schema.Set).List()
	documentTypes := []int{}
	for _, tag := range dTypes {
		t := tag.(int)
		documentTypes = append(documentTypes, t)
	}

	for _, docType := range documentTypes {
		c.AddWorkflowIndexDocumentType(workflowTemplate.ID, docType)
	}

	return resourceWorkflowTemplateRead(d, m)
}

func resourceWorkflowTemplateRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetWorkflowTemplateById(id)
	if err != nil {
		return err
	}

	docTypes, err := c.GetWorkflowIndexDocumentTypes(source.ID)
	if err != nil {
		return err
	}

	if err := d.Set("document_types", docTypes); err != nil {
		return err
	}

	return workflowTemplateToData(source, d)
}

func resourceWorkflowTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplate := dataToWorkflowTemplate(d)
	workflowTemplate.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateWorkflowTemplate(*workflowTemplate)
	if err != nil {
		return err
	}
	if d.HasChange("document_types") {
		o, n := d.GetChange("document_types")

		oTypes := o.(*schema.Set)
		nTypes := n.(*schema.Set)

		removals := oTypes.Difference(nTypes)
		for _, removal := range removals.List() {
			if err := c.RemoveWorkflowIndexDocumentType(workflowTemplate.ID, removal.(int)); err != nil {
				return err
			}
		}

		additions := nTypes.Difference(oTypes)
		for _, addition := range additions.List() {
			if err := c.AddWorkflowIndexDocumentType(workflowTemplate.ID, addition.(int)); err != nil {
				return err
			}
		}
	}

	return resourceWorkflowTemplateRead(d, m)
}

func resourceWorkflowTemplateDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteWorkflowTemplate(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceWorkflowTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	workflowTemplate, err := c.GetWorkflowTemplateById(id)
	if err != nil {
		return rd, err
	}

	docTypes, err := c.GetWorkflowIndexDocumentTypes(workflowTemplate.ID)
	if err != nil {
		return rd, err
	}

	if err := d.Set("document_types", docTypes); err != nil {
		return rd, err
	}

	err = workflowTemplateToData(workflowTemplate, d)
	return rd, err
}

func workflowTemplateToData(workflowTemplate *client.WorkflowTemplate, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", workflowTemplate.ID))
	if err := d.Set("label", workflowTemplate.Label); err != nil {
		return err
	}
	if err := d.Set("internal_name", workflowTemplate.InternalName); err != nil {
		return err
	}

	return nil
}

func dataToWorkflowTemplate(d *schema.ResourceData) *client.WorkflowTemplate {
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.WorkflowTemplate{
		ID:           id,
		Label:        d.Get("label").(string),
		InternalName: d.Get("internal_name").(string),
	}

	return &newDocType
}
