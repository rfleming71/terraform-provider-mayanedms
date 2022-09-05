package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type MayanEdmsClient interface {
	GetDocumentTypeById(id int) (*DocumentType, error)
	CreateDocumentType(documentType DocumentType) (*DocumentType, error)
	UpdateDocumentType(documentType DocumentType) (*DocumentType, error)
	DeleteDocumentType(id int) error

	GetSourceById(id int) (*Source, error)
	CreateSource(source Source) (*Source, error)
	UpdateSource(documentType Source) (*Source, error)
	DeleteSource(id int) error

	GetTagById(id int) (*Tag, error)
	CreateTag(tag Tag) (*Tag, error)
	UpdateTag(tag Tag) (*Tag, error)
	DeleteTag(id int) error

	GetIndexTemplateById(id int) (*IndexTemplate, error)
	CreateIndexTemplate(indexTemplate IndexTemplate) (*IndexTemplate, error)
	UpdateIndexTemplate(indexTemplate IndexTemplate) (*IndexTemplate, error)
	DeleteIndexTemplate(id int) error

	GetIndexTemplateDocumentTypes(indexTemplateId int) ([]int, error)
	AddIndexTemplateDocumentType(indexTemplateId int, documentTypeId int) error
	RemoveIndexTemplateDocumentType(indexTemplateId int, documentTypeId int) error

	GetIndexTemplateNodeById(indexId, nodeId int) (*IndexTemplateNode, error)
	CreateIndexTemplateNode(indexTemplateNode IndexTemplateNode) (*IndexTemplateNode, error)
	UpdateIndexTemplateNode(indexTemplateNode IndexTemplateNode) (*IndexTemplateNode, error)
	DeleteIndexTemplateNode(indexId, nodeId int) error

	GetGroupById(id int) (*Group, error)
	CreateGroup(group Group) (*Group, error)
	UpdateGroup(group Group) (*Group, error)
	DeleteGroup(id int) error
	GetGroupUsers(groupId int) ([]int, error)
	AddGroupUser(groupId int, userId int) error
	RemoveGroupUser(groupId int, userId int) error

	GetWorkflowTemplateById(id int) (*WorkflowTemplate, error)
	CreateWorkflowTemplate(workflowTemplate WorkflowTemplate) (*WorkflowTemplate, error)
	UpdateWorkflowTemplate(workflowTemplate WorkflowTemplate) (*WorkflowTemplate, error)
	DeleteWorkflowTemplate(id int) error
	GetWorkflowIndexDocumentTypes(workflowTemplateId int) ([]int, error)
	AddWorkflowIndexDocumentType(workflowTemplateId int, documentTypeId int) error
	RemoveWorkflowIndexDocumentType(workflowTemplateId int, documentTypeId int) error

	GetWorkflowTemplateState(workflowTemplateId int, stateId int) (*WorkflowTemplateState, error)
	CreateWorkflowTemplateState(workflowTemplateId int, state WorkflowTemplateState) (*WorkflowTemplateState, error)
	RemoveWorkflowTemplateState(workflowTemplateId int, stateId int) error
	UpdateWorkflowTemplateState(workflowTemplateId int, state WorkflowTemplateState) (*WorkflowTemplateState, error)

	GetWorkflowTemplateTransition(workflowTemplateId int, transitionId int) (*WorkflowTemplateTransition, error)
	CreateWorkflowTemplateTransition(workflowTemplateId int, transition WorkflowTemplateTransition) (*WorkflowTemplateTransition, error)
	RemoveWorkflowTemplateTransition(workflowTemplateId int, transitionId int) error
	UpdateWorkflowTemplateTransition(workflowTemplateId int, transition WorkflowTemplateTransition) (*WorkflowTemplateTransition, error)

	GetRoleById(id int) (*Role, error)
	CreateRole(tag Role) (*Role, error)
	UpdateRole(tag Role) (*Role, error)
	DeleteRole(id int) error
	GetRoleGroups(roleId int) ([]int, error)
	AddRoleGroup(roleId int, groupId int) error
	RemoveRoleGroup(roleId int, groupId int) error
	GetRolePermissions(roleId int) ([]string, error)
	AddRolePermission(roleId int, permissionPk string) error
	RemoveRolePermission(roleId int, permissionPk string) error

	CreateMetadataType(metadataType MetadataType) (*MetadataType, error)
	GetMetadataTypeById(id int) (*MetadataType, error)
	DeleteMetadataType(id int) error
	UpdateMetadataType(metadataType MetadataType) (*MetadataType, error)
}

type ClientConfig struct {
	Url                string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

type Client struct {
	client *http.Client
	url    string
	token  string
}

func NewMayanEdmsClient(config ClientConfig) (MayanEdmsClient, error) {
	client := &Client{
		client: &http.Client{},
		url:    config.Url + "/api/v4/",
	}

	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: config.Username,
		Password: config.Password,
	}

	response := struct {
		Token string `json:"token"`
	}{}

	err := client.performRequest("auth/token/obtain/", http.MethodPost, &request, &response)
	if err != nil {
		return &Client{}, fmt.Errorf("failed to obtain token: %v", err)
	}

	client.token = response.Token

	return client, nil
}

func (c *Client) performRequest(path string, method string, body interface{}, response interface{}) error {

	var req *http.Request
	var err error
	if body != nil {
		json_data, err := json.Marshal(body)

		if err != nil {
			return err
		}

		req, err = http.NewRequest(method, c.url+path, bytes.NewBuffer(json_data))

		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequest(method, c.url+path, nil)
		if err != nil {
			return err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	if c.token != "" {
		req.Header.Add("Authorization", "Token "+c.token)
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)
	}
	return err
}
