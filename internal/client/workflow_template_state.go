package client

import (
	"fmt"
	"net/http"
)

type WorkflowTemplateState struct {
	ID         int    `json:"id"`
	Completion int    `json:"completion"`
	Initial    bool   `json:"initial"`
	Label      string `json:"label"`
}

func (c *Client) GetWorkflowTemplateState(workflowTemplateId int, stateId int) (*WorkflowTemplateState, error) {
	var result WorkflowTemplateState
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/states/%v/", workflowTemplateId, stateId), http.MethodGet, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateWorkflowTemplateState(workflowTemplateId int, state WorkflowTemplateState) (*WorkflowTemplateState, error) {
	var newState WorkflowTemplateState
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/states/", workflowTemplateId), http.MethodPost, &state, &newState)
	if err != nil {
		return &WorkflowTemplateState{}, err
	}

	return &newState, nil
}

func (c *Client) RemoveWorkflowTemplateState(workflowTemplateId int, stateId int) error {

	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/states/%v/", workflowTemplateId, stateId), http.MethodDelete, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateWorkflowTemplateState(workflowTemplateId int, state WorkflowTemplateState) (*WorkflowTemplateState, error) {
	var updatedState WorkflowTemplateState
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/states/%v/", workflowTemplateId, state.ID), http.MethodPut, &state, &updatedState)
	if err != nil {
		return &WorkflowTemplateState{}, err
	}

	return &updatedState, nil
}
