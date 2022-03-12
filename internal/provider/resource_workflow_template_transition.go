package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceWorkflowTemplateTransition() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkflowTemplateTransitionCreate,
		Read:   resourceWorkflowTemplateTransitionRead,
		Update: resourceWorkflowTemplateTransitionUpdate,
		Delete: resourceWorkflowTemplateTransitionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWorkflowTemplateTransitionImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_state": {
				Type:     schema.TypeString,
				Required: true,
			},
			"origin_state": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"workflow_template": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWorkflowTemplateTransitionCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, newWorkflowTemplateTransition := dataToWorkflowTemplateTransition(d)

	workflowTemplateTransition, err := c.CreateWorkflowTemplateTransition(workflowTemplateId, *newWorkflowTemplateTransition)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v-%v", workflowTemplateId, workflowTemplateTransition.ID))

	return resourceWorkflowTemplateTransitionRead(d, m)
}

func resourceWorkflowTemplateTransitionRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, stateId, err := getIdInformation(d)
	if err != nil {
		return err
	}

	source, err := c.GetWorkflowTemplateTransition(workflowTemplateId, stateId)
	if err != nil {
		return err
	}

	return workflowTemplateTransitionToData(workflowTemplateId, source, d)
}

func resourceWorkflowTemplateTransitionUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	_, workflowTemplateTransition := dataToWorkflowTemplateTransition(d)
	workflowTemplateId, stateId, err := getIdInformation(d)
	if err != nil {
		return err
	}
	workflowTemplateTransition.ID = stateId

	_, err = c.UpdateWorkflowTemplateTransition(workflowTemplateId, *workflowTemplateTransition)
	if err != nil {
		return err
	}

	return resourceWorkflowTemplateTransitionRead(d, m)
}

func resourceWorkflowTemplateTransitionDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, stateId, err := getIdInformation(d)
	if err != nil {
		return err
	}

	err = c.RemoveWorkflowTemplateTransition(workflowTemplateId, stateId)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceWorkflowTemplateTransitionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, stateId, err := getIdInformation(d)
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	workflowTemplateTransition, err := c.GetWorkflowTemplateTransition(workflowTemplateId, stateId)
	if err != nil {
		return rd, err
	}

	err = workflowTemplateTransitionToData(workflowTemplateId, workflowTemplateTransition, d)
	return rd, err
}

func workflowTemplateTransitionToData(workflowTemplateId int, workflowTemplateTransition *client.WorkflowTemplateTransition, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v-%v", workflowTemplateId, workflowTemplateTransition.ID))
	if err := d.Set("label", workflowTemplateTransition.Label); err != nil {
		return err
	}

	tmp := fmt.Sprintf("%v-%v", workflowTemplateId, workflowTemplateTransition.DestinationState.ID)
	if err := d.Set("destination_state", tmp); err != nil {
		return err
	}

	tmp = fmt.Sprintf("%v-%v", workflowTemplateId, workflowTemplateTransition.OriginState.ID)
	if err := d.Set("origin_state", tmp); err != nil {
		return err
	}

	if err := d.Set("condition", workflowTemplateTransition.Condition); err != nil {
		return err
	}

	if err := d.Set("workflow_template", workflowTemplateId); err != nil {
		return err
	}

	return nil
}

func dataToWorkflowTemplateTransition(d *schema.ResourceData) (int, *client.WorkflowTemplateTransition) {
	_, id, _ := getIdInformation(d)
	newDocType := client.WorkflowTemplateTransition{
		ID:        id,
		Label:     d.Get("label").(string),
		Condition: d.Get("condition").(string),
	}

	_, id, _ = breakCompositeId(d.Get("destination_state").(string))
	newDocType.DestinationState = client.WorkflowTemplateState{
		ID: id,
	}

	_, id, _ = breakCompositeId(d.Get("origin_state").(string))
	newDocType.OriginState = client.WorkflowTemplateState{
		ID: id,
	}

	return d.Get("workflow_template").(int), &newDocType
}
