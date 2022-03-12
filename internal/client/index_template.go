package client

import (
	"fmt"
	"net/http"
)

type IndexTemplate struct {
	ID      int    `json:"id"`
	Enabled bool   `json:"enabled"`
	Label   string `json:"label"`
	Slug    string `json:"slug"`
}

type IndexTemplateNode struct {
	ID            int                  `json:"id"`
	Expression    string               `json:"expression"`
	LinkDocuments bool                 `json:"link_documents"`
	Children      []*IndexTemplateNode `json:"children"`
}

func (c *Client) CreateIndexTemplate(tag IndexTemplate) (*IndexTemplate, error) {
	var createdIndex *IndexTemplate
	err := c.performRequest("index_templates/", http.MethodPost, &tag, &createdIndex)
	if err != nil {
		return &IndexTemplate{}, err
	}

	return createdIndex, nil
}

func (c *Client) GetIndexTemplateById(id int) (*IndexTemplate, error) {
	var tag *IndexTemplate
	err := c.performRequest(fmt.Sprintf("index_templates/%v/", id), http.MethodGet, nil, &tag)
	if err != nil {
		return &IndexTemplate{}, err
	}

	return tag, nil
}

func (c *Client) DeleteIndexTemplate(id int) error {
	err := c.performRequest(fmt.Sprintf("index_templates/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateIndexTemplate(documentType IndexTemplate) (*IndexTemplate, error) {
	var updatedIndexTemplate *IndexTemplate
	err := c.performRequest(fmt.Sprintf("index_templates/%v/", documentType.ID), http.MethodPut, &documentType, &updatedIndexTemplate)
	if err != nil {
		return &IndexTemplate{}, err
	}

	return updatedIndexTemplate, nil
}

func (c *Client) GetIndexTemplateDocumentTypes(indexTemplateId int) ([]int, error) {
	var results struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	err := c.performRequest(fmt.Sprintf("index_templates/%v/document_types/", indexTemplateId), http.MethodGet, nil, &results)
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, result := range results.Results {
		ids = append(ids, result.ID)
	}

	return ids, nil
}

func (c *Client) AddIndexTemplateDocumentType(indexTemplateId int, documentTypeId int) error {

	var request struct {
		DocumentTypeId int `json:"document_type"`
	}
	request.DocumentTypeId = documentTypeId

	err := c.performRequest(fmt.Sprintf("index_templates/%v/document_types/add/", indexTemplateId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveIndexTemplateDocumentType(indexTemplateId int, documentTypeId int) error {

	var request struct {
		DocumentTypeId int `json:"document_type"`
	}
	request.DocumentTypeId = documentTypeId

	err := c.performRequest(fmt.Sprintf("index_templates/%v/document_types/remove/", indexTemplateId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}
