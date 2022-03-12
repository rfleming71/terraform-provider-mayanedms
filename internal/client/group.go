package client

import (
	"fmt"
	"net/http"
)

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) CreateGroup(group Group) (*Group, error) {
	var createdGroup *Group
	err := c.performRequest("groups/", http.MethodPost, &group, &createdGroup)
	if err != nil {
		return &Group{}, err
	}

	return createdGroup, nil
}

func (c *Client) GetGroupById(id int) (*Group, error) {
	var group *Group
	err := c.performRequest(fmt.Sprintf("groups/%v/", id), http.MethodGet, nil, &group)
	if err != nil {
		return &Group{}, err
	}

	return group, nil
}

func (c *Client) DeleteGroup(id int) error {
	err := c.performRequest(fmt.Sprintf("groups/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateGroup(group Group) (*Group, error) {
	var updatedGroup *Group
	err := c.performRequest(fmt.Sprintf("groups/%v/", group.ID), http.MethodPut, &group, &updatedGroup)
	if err != nil {
		return &Group{}, err
	}

	return updatedGroup, nil
}

func (c *Client) GetGroupUsers(groupId int) ([]int, error) {
	var results struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	err := c.performRequest(fmt.Sprintf("groups/%v/users/", groupId), http.MethodGet, nil, &results)
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, result := range results.Results {
		ids = append(ids, result.ID)
	}

	return ids, nil
}

func (c *Client) AddGroupUser(groupId int, userId int) error {

	var request struct {
		UserId int `json:"user"`
	}
	request.UserId = userId

	err := c.performRequest(fmt.Sprintf("groups/%v/users/add/", groupId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveGroupUser(groupId int, userId int) error {

	var request struct {
		UserId int `json:"user"`
	}
	request.UserId = userId

	err := c.performRequest(fmt.Sprintf("groups/%v/users/remove/", groupId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}
