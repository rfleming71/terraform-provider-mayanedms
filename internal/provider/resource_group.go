package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newGroup := dataToGroup(d)

	group, err := c.CreateGroup(*newGroup)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", group.ID))

	ids := d.Get("users").(*schema.Set).List()
	userIds := []int{}
	for _, id := range ids {
		t := id.(int)
		userIds = append(userIds, t)
	}

	for _, id := range userIds {
		c.AddGroupUser(group.ID, id)
	}

	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	group, err := c.GetGroupById(id)
	if err != nil {
		return err
	}

	userIds, err := c.GetGroupUsers(group.ID)
	if err != nil {
		return err
	}

	if err := d.Set("users", userIds); err != nil {
		return err
	}

	return groupToData(group, d)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	group := dataToGroup(d)
	group.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateGroup(*group)
	if err != nil {
		return err
	}

	if d.HasChange("users") {
		o, n := d.GetChange("users")

		oTypes := o.(*schema.Set)
		nTypes := n.(*schema.Set)

		removals := oTypes.Difference(nTypes)
		for _, removal := range removals.List() {
			if err := c.RemoveGroupUser(group.ID, removal.(int)); err != nil {
				return err
			}
		}

		additions := nTypes.Difference(oTypes)
		for _, addition := range additions.List() {
			if err := c.AddGroupUser(group.ID, addition.(int)); err != nil {
				return err
			}
		}
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteGroup(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	group, err := c.GetGroupById(id)
	if err != nil {
		return rd, err
	}

	userIds, err := c.GetGroupUsers(group.ID)
	if err != nil {
		return rd, err
	}

	if err := d.Set("users", userIds); err != nil {
		return rd, err
	}

	err = groupToData(group, d)
	return rd, err
}

func groupToData(group *client.Group, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", group.ID))
	if err := d.Set("name", group.Name); err != nil {
		return err
	}

	return nil
}

func dataToGroup(d *schema.ResourceData) *client.Group {
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.Group{
		ID:   id,
		Name: d.Get("name").(string),
	}

	return &newDocType
}
