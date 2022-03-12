package client

import (
	"fmt"
	"net/http"
)

type MetadataType struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Name      string `json:"name"`
	Default   string `json:"default"`
	Lookup    string `json:"lookup"`
	Validator string `json:"validator"`
	Parser    string `json:"parser"`
}

func (c *Client) CreateMetadataType(metadataType MetadataType) (*MetadataType, error) {
	var createdType *MetadataType
	err := c.performRequest("metadata_types/", http.MethodPost, &metadataType, &createdType)
	if err != nil {
		return &MetadataType{}, err
	}

	return createdType, nil
}

func (c *Client) GetMetadataTypeById(id int) (*MetadataType, error) {
	var metadataType *MetadataType
	err := c.performRequest(fmt.Sprintf("metadata_types/%v/", id), http.MethodGet, nil, &metadataType)
	if err != nil {
		return &MetadataType{}, err
	}

	return metadataType, nil
}

func (c *Client) DeleteMetadataType(id int) error {
	err := c.performRequest(fmt.Sprintf("metadata_types/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateMetadataType(metadataType MetadataType) (*MetadataType, error) {
	var updatedMetadataType *MetadataType
	err := c.performRequest(fmt.Sprintf("metadata_types/%v/", metadataType.ID), http.MethodPut, &metadataType, &updatedMetadataType)
	if err != nil {
		return &MetadataType{}, err
	}

	return updatedMetadataType, nil
}
