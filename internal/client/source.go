package client

import (
	"fmt"
	"net/http"
)

type Source struct {
	ID          int    `json:"id"`
	Label       string `json:"label"`
	BackendData string `json:"backend_data"`
	BackendPath string `json:"backend_path"`
	Enabled     bool   `json:"enabled"`
}

func (c *Client) CreateSource(source Source) (*Source, error) {
	var createdSource *Source
	err := c.performRequest("sources/", http.MethodPost, &source, &createdSource)
	if err != nil {
		return &Source{}, err
	}

	return createdSource, nil
}

func (c *Client) GetSourceById(id int) (*Source, error) {
	var source *Source
	err := c.performRequest(fmt.Sprintf("sources/%v/", id), http.MethodGet, nil, &source)
	if err != nil {
		return &Source{}, err
	}

	return source, nil
}

func (c *Client) DeleteSource(id int) error {
	err := c.performRequest(fmt.Sprintf("sources/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateSource(source Source) (*Source, error) {
	var updatedSource *Source
	err := c.performRequest(fmt.Sprintf("sources/%v/", source.ID), http.MethodPut, &source, &updatedSource)
	if err != nil {
		return &Source{}, err
	}

	return updatedSource, nil
}
