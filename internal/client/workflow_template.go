package client

import (
	"fmt"
	"net/http"
)

type WorkflowTemplate struct {
	ID           int    `json:"id"`
	Label        string `json:"label"`
	InternalName string `json:"internal_name"`
}

type WorkflowStateAction struct {
	ID         int    `json:"id"`
	ActionPath string `json:"action_path"`
	ActionData string `json:"action_data"`
	Enabled    bool   `json:"enabled"`
	Label      string `json:"label"`
	When       int    `json:"when"`
}

func (c *Client) CreateWorkflowTemplate(workflowTemplate WorkflowTemplate) (*WorkflowTemplate, error) {
	var createdDoc *WorkflowTemplate
	err := c.performRequest("workflow_templates/", http.MethodPost, &workflowTemplate, &createdDoc)
	if err != nil {
		return &WorkflowTemplate{}, err
	}

	return createdDoc, nil
}

func (c *Client) GetWorkflowTemplateById(id int) (*WorkflowTemplate, error) {
	var workflowTemplate *WorkflowTemplate
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/", id), http.MethodGet, nil, &workflowTemplate)
	if err != nil {
		return &WorkflowTemplate{}, err
	}

	return workflowTemplate, nil
}

func (c *Client) DeleteWorkflowTemplate(id int) error {
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateWorkflowTemplate(documentType WorkflowTemplate) (*WorkflowTemplate, error) {
	var updatedWorkflowTemplate *WorkflowTemplate
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/", documentType.ID), http.MethodPut, &documentType, &updatedWorkflowTemplate)
	if err != nil {
		return &WorkflowTemplate{}, err
	}

	return updatedWorkflowTemplate, nil
}

func (c *Client) GetWorkflowIndexDocumentTypes(workflowTemplateId int) ([]int, error) {
	var results struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/document_types/", workflowTemplateId), http.MethodGet, nil, &results)
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, result := range results.Results {
		ids = append(ids, result.ID)
	}

	return ids, nil
}

func (c *Client) AddWorkflowIndexDocumentType(workflowTemplateId int, documentTypeId int) error {

	var request struct {
		DocumentTypeId int `json:"document_type_id"`
	}
	request.DocumentTypeId = documentTypeId

	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/document_types/add/", workflowTemplateId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveWorkflowIndexDocumentType(workflowTemplateId int, documentTypeId int) error {

	var request struct {
		DocumentTypeId int `json:"document_type_id"`
	}
	request.DocumentTypeId = documentTypeId

	err := c.performRequest(fmt.Sprintf("workflow_templates/%v/document_types/remove/", workflowTemplateId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}
