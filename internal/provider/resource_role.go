package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoleImport,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	newRole := dataToRole(d)

	role, err := c.CreateRole(*newRole)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", role.ID))
	groups := d.Get("groups").(*schema.Set).List()
	for _, group := range groups {
		c.AddRoleGroup(newRole.ID, group.(int))
	}

	permissions := d.Get("permissions").(*schema.Set).List()
	for _, permission := range permissions {
		c.AddRolePermission(newRole.ID, permission.(string))
	}

	return resourceRoleRead(d, m)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())

	source, err := c.GetRoleById(id)
	if err != nil {
		return err
	}

	groups, err := c.GetRoleGroups(source.ID)
	if err != nil {
		return err
	}

	if err := d.Set("groups", groups); err != nil {
		return err
	}

	permissions, err := c.GetRolePermissions(source.ID)
	if err != nil {
		return err
	}

	if err := d.Set("permissions", permissions); err != nil {
		return err
	}

	return roleToData(source, d)
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	role := dataToRole(d)
	role.ID, _ = strconv.Atoi(d.Id())
	_, err := c.UpdateRole(*role)
	if err != nil {
		return err
	}

	if d.HasChange("groups") {
		o, n := d.GetChange("groups")

		oTypes := o.(*schema.Set)
		nTypes := n.(*schema.Set)

		removals := oTypes.Difference(nTypes)
		for _, removal := range removals.List() {
			if err := c.RemoveRoleGroup(role.ID, removal.(int)); err != nil {
				return err
			}
		}

		additions := nTypes.Difference(oTypes)
		for _, addition := range additions.List() {
			if err := c.AddRoleGroup(role.ID, addition.(int)); err != nil {
				return err
			}
		}
	}

	if d.HasChange("permissions") {
		o, n := d.GetChange("permissions")

		oTypes := o.(*schema.Set)
		nTypes := n.(*schema.Set)

		removals := oTypes.Difference(nTypes)
		for _, removal := range removals.List() {
			if err := c.RemoveRolePermission(role.ID, removal.(string)); err != nil {
				return err
			}
		}

		additions := nTypes.Difference(oTypes)
		for _, addition := range additions.List() {
			if err := c.AddRolePermission(role.ID, addition.(string)); err != nil {
				return err
			}
		}
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)
	id, _ := strconv.Atoi(d.Id())
	err := c.DeleteRole(id)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)
	id, err := strconv.Atoi(d.Id())
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	role, err := c.GetRoleById(id)
	if err != nil {
		return rd, err
	}

	groups, err := c.GetRoleGroups(role.ID)
	if err != nil {
		return rd, err
	}

	if err := d.Set("groups", groups); err != nil {
		return rd, err
	}

	permissions, err := c.GetRolePermissions(role.ID)
	if err != nil {
		return rd, err
	}

	if err := d.Set("permissions", permissions); err != nil {
		return rd, err
	}

	err = roleToData(role, d)
	return rd, err
}

func roleToData(role *client.Role, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v", role.ID))
	if err := d.Set("label", role.Label); err != nil {
		return err
	}

	return nil
}

func dataToRole(d *schema.ResourceData) *client.Role {
	id, _ := strconv.Atoi(d.Id())
	newDocType := client.Role{
		ID:    id,
		Label: d.Get("label").(string),
	}

	return &newDocType
}
