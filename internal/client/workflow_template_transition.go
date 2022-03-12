package client

import (
	"fmt"
	"net/http"
)

type WorkflowTemplateTransition struct {
	ID               int                   `json:"id"`
	Condition        string                `json:"condition"`
	DestinationState WorkflowTemplateState `json:"destination_state"`
	Label            string                `json:"label"`
	OriginState      WorkflowTemplateState `json:"origin_state"`
}

type workflowTemplateTransition struct {
	ID                 int    `json:"id"`
	Condition          string `json:"condition"`
	DestinationStateId int    `json:"destination_state_id"`
	Label              string `json:"label"`
	OriginStateId      int    `json:"origin_state_id"`
}

func (c *Client) GetWorkflowTemplateTransition(workflowTemplateId int, transitionId int) (*WorkflowTemplateTransition, error) {
	var result WorkflowTemplateTransition
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/transitions/%v/", workflowTemplateId, transitionId), http.MethodGet, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateWorkflowTemplateTransition(workflowTemplateId int, transition WorkflowTemplateTransition) (*WorkflowTemplateTransition, error) {
	var newTransition WorkflowTemplateTransition
	request := workflowTemplateTransition{
		Label:              transition.Label,
		Condition:          transition.Condition,
		DestinationStateId: transition.DestinationState.ID,
		OriginStateId:      transition.OriginState.ID,
	}
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/transitions/", workflowTemplateId), http.MethodPost, &request, &newTransition)
	if err != nil {
		return &WorkflowTemplateTransition{}, err
	}

	return &newTransition, nil
}

func (c *Client) RemoveWorkflowTemplateTransition(workflowTemplateId int, transitionId int) error {

	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/transitions/%v/", workflowTemplateId, transitionId), http.MethodDelete, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateWorkflowTemplateTransition(workflowTemplateId int, transition WorkflowTemplateTransition) (*WorkflowTemplateTransition, error) {
	var updatedTransition WorkflowTemplateTransition
	request := workflowTemplateTransition{
		Label:              transition.Label,
		Condition:          transition.Condition,
		DestinationStateId: transition.DestinationState.ID,
		OriginStateId:      transition.OriginState.ID,
	}
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/transitions/%v/", workflowTemplateId, transition.ID), http.MethodPut, &request, &updatedTransition)
	if err != nil {
		return &WorkflowTemplateTransition{}, err
	}

	return &updatedTransition, nil
}
