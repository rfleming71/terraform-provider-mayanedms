package client

import (
	"fmt"
	"net/http"
)

type Role struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func (c *Client) CreateRole(role Role) (*Role, error) {
	var createdDoc *Role
	err := c.performRequest("roles/", http.MethodPost, &role, &createdDoc)
	if err != nil {
		return &Role{}, err
	}

	return createdDoc, nil
}

func (c *Client) GetRoleById(id int) (*Role, error) {
	var role *Role
	err := c.performRequest(fmt.Sprintf("roles/%v/", id), http.MethodGet, nil, &role)
	if err != nil {
		return &Role{}, err
	}

	return role, nil
}

func (c *Client) DeleteRole(id int) error {
	err := c.performRequest(fmt.Sprintf("roles/%v/", id), http.MethodDelete, nil, nil)
	return err
}

func (c *Client) UpdateRole(documentType Role) (*Role, error) {
	var updatedRole *Role
	err := c.performRequest(fmt.Sprintf("roles/%v/", documentType.ID), http.MethodPut, &documentType, &updatedRole)
	if err != nil {
		return &Role{}, err
	}

	return updatedRole, nil
}

func (c *Client) GetRoleGroups(roleId int) ([]int, error) {
	var results struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	err := c.performRequest(fmt.Sprintf("roles/%v/groups/?page_size=200", roleId), http.MethodGet, nil, &results)
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, result := range results.Results {
		ids = append(ids, result.ID)
	}

	return ids, nil
}

func (c *Client) AddRoleGroup(roleId int, groupId int) error {

	var request struct {
		GroupId int `json:"group_id"`
	}
	request.GroupId = groupId

	err := c.performRequest(fmt.Sprintf("roles/%v/groups/add/", roleId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveRoleGroup(roleId int, groupId int) error {

	var request struct {
		GroupId int `json:"group_id"`
	}
	request.GroupId = groupId

	err := c.performRequest(fmt.Sprintf("roles/%v/groups/remove/", roleId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetRolePermissions(roleId int) ([]string, error) {
	var results struct {
		Results []struct {
			Pk string `json:"pk"`
		} `json:"results"`
	}
	err := c.performRequest(fmt.Sprintf("roles/%v/permissions/?page_size=200", roleId), http.MethodGet, nil, &results)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, result := range results.Results {
		ids = append(ids, result.Pk)
	}

	return ids, nil
}

func (c *Client) AddRolePermission(roleId int, permissionPk string) error {

	var request struct {
		Permission string `json:"permission"`
	}
	request.Permission = permissionPk

	err := c.performRequest(fmt.Sprintf("roles/%v/permissions/add/", roleId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveRolePermission(roleId int, permissionPk string) error {

	var request struct {
		Permission string `json:"permission"`
	}
	request.Permission = permissionPk

	err := c.performRequest(fmt.Sprintf("roles/%v/permissions/remove/", roleId), http.MethodPost, &request, nil)
	if err != nil {
		return err
	}

	return nil
}
