package client

import (
	"fmt"
	"net/http"
)

type DocumentType struct {
	ID                                int     `json:"id"`
	DeleteTimePeriod                  int     `json:"delete_time_period"`
	DeleteTimeUnit                    string  `json:"delete_time_unit"`
	Label                             string  `json:"label"`
	TrashTimePeriod                   *int    `json:"trash_time_period"`
	TrashTimeUnit                     *string `json:"trash_time_unit"`
	FileNameGeneratorBackend          string  `json:"filename_generator_backend"`
	FileNameGeneratorBackendArguments string  `json:"filename_generator_backend_arguments"`
}

func (c *Client) CreateDocumentType(documentType DocumentType) (*DocumentType, error) {
	var createdDoc *DocumentType
	err := c.performRequest("document_types/", http.MethodPost, &documentType, &createdDoc)
	if err != nil {
		return &DocumentType{}, err
	}

	return createdDoc, nil
}

func (c *Client) GetDocumentTypeById(id int) (*DocumentType, error) {
	var documentType *DocumentType
	err := c.performRequest(fmt.Sprintf("document_types/%v/", id), http.MethodGet, nil, &documentType)
	if err != nil {
		return &DocumentType{}, err
	}

	return documentType, nil
}

func (c *Client) DeleteDocumentType(id int) error {
	err := c.performRequest(fmt.Sprintf("document_types/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateDocumentType(documentType DocumentType) (*DocumentType, error) {
	var updatedDocType *DocumentType
	err := c.performRequest(fmt.Sprintf("document_types/%v/", documentType.ID), http.MethodPut, &documentType, &updatedDocType)
	if err != nil {
		return &DocumentType{}, err
	}

	return updatedDocType, nil
}
