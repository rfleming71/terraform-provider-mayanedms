package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceWorkflowTemplateState() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkflowTemplateStateCreate,
		Read:   resourceWorkflowTemplateStateRead,
		Update: resourceWorkflowTemplateStateUpdate,
		Delete: resourceWorkflowTemplateStateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWorkflowTemplateStateImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"completion": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"initial": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"workflow_template": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWorkflowTemplateStateCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, newWorkflowTemplateState := dataToWorkflowTemplateState(d)

	workflowTemplateState, err := c.CreateWorkflowTemplateState(workflowTemplateId, *newWorkflowTemplateState)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v-%v", workflowTemplateId, workflowTemplateState.ID))

	return resourceWorkflowTemplateStateRead(d, m)
}

func resourceWorkflowTemplateStateRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, stateId, err := getIdInformation(d)
	if err != nil {
		return err
	}

	source, err := c.GetWorkflowTemplateState(workflowTemplateId, stateId)
	if err != nil {
		return err
	}

	return workflowTemplateStateToData(workflowTemplateId, source, d)
}

func resourceWorkflowTemplateStateUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	_, workflowTemplateState := dataToWorkflowTemplateState(d)
	workflowTemplateId, stateId, err := getIdInformation(d)
	if err != nil {
		return err
	}
	workflowTemplateState.ID = stateId

	_, err = c.UpdateWorkflowTemplateState(workflowTemplateId, *workflowTemplateState)
	if err != nil {
		return err
	}

	return resourceWorkflowTemplateStateRead(d, m)
}

func resourceWorkflowTemplateStateDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, stateId, err := getIdInformation(d)
	if err != nil {
		return err
	}

	err = c.RemoveWorkflowTemplateState(workflowTemplateId, stateId)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceWorkflowTemplateStateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	workflowTemplateId, stateId, err := getIdInformation(d)
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	workflowTemplateState, err := c.GetWorkflowTemplateState(workflowTemplateId, stateId)
	if err != nil {
		return rd, err
	}

	err = workflowTemplateStateToData(workflowTemplateId, workflowTemplateState, d)
	return rd, err
}

func workflowTemplateStateToData(workflowTemplateId int, workflowTemplateState *client.WorkflowTemplateState, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v-%v", workflowTemplateId, workflowTemplateState.ID))
	if err := d.Set("label", workflowTemplateState.Label); err != nil {
		return err
	}

	if err := d.Set("initial", workflowTemplateState.Initial); err != nil {
		return err
	}

	if err := d.Set("completion", workflowTemplateState.Completion); err != nil {
		return err
	}

	if err := d.Set("workflow_template", workflowTemplateId); err != nil {
		return err
	}

	return nil
}

func dataToWorkflowTemplateState(d *schema.ResourceData) (int, *client.WorkflowTemplateState) {
	_, id, _ := getIdInformation(d)
	newDocType := client.WorkflowTemplateState{
		ID:         id,
		Label:      d.Get("label").(string),
		Initial:    d.Get("initial").(bool),
		Completion: d.Get("completion").(int),
	}

	return d.Get("workflow_template").(int), &newDocType
}

func getIdInformation(d *schema.ResourceData) (int, int, error) {
	return breakCompositeId(d.Id())
}

func breakCompositeId(id string) (int, int, error) {
	ids := strings.Split(id, "-")
	part1, err := strconv.Atoi(ids[0])
	if err != nil {
		return 0, 0, err
	}

	part2, err := strconv.Atoi(ids[1])
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil

}
