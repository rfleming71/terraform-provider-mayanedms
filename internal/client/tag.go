package client

import (
	"fmt"
	"net/http"
)

type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
}

func (c *Client) CreateTag(tag Tag) (*Tag, error) {
	var createdDoc *Tag
	err := c.performRequest("tags/", http.MethodPost, &tag, &createdDoc)
	if err != nil {
		return &Tag{}, err
	}

	return createdDoc, nil
}

func (c *Client) GetTagById(id int) (*Tag, error) {
	var tag *Tag
	err := c.performRequest(fmt.Sprintf("tags/%v/", id), http.MethodGet, nil, &tag)
	if err != nil {
		return &Tag{}, err
	}

	return tag, nil
}

func (c *Client) DeleteTag(id int) error {
	err := c.performRequest(fmt.Sprintf("tags/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateTag(documentType Tag) (*Tag, error) {
	var updatedTag *Tag
	err := c.performRequest(fmt.Sprintf("tags/%v/", documentType.ID), http.MethodPut, &documentType, &updatedTag)
	if err != nil {
		return &Tag{}, err
	}

	return updatedTag, nil
}
