package controllers

import "fmt"

// @provider({"id":"main.role.name"})
func NewRoleName() string {
	return "role"
}

type RoleController struct {
	Name string
}

// @group({"id":"group.controllers"})
// @inject({"param":"name", "id":"main.role.name"})
func NewRoleController(name string) *RoleController {
	return &RoleController{
		Name: name,
	}
}

func (c *RoleController) Print() {
	fmt.Printf("role controller name:%s\n", c.Name)
}
