package main

import (
	"fmt"

	"github.com/werbenhu/digo"
)

type Controller interface {
	Print()
}

// @provider({"id":"main.user.name"})
func NewUserName() string {
	return "user"
}

type UserController struct {
	Name string
}

// @group({"id":"group.controllers"})
// @inject({"param":"name", "id":"main.user.name"})
func NewUserController(name string) *UserController {
	return &UserController{
		Name: name,
	}
}

func (c *UserController) Print() {
	fmt.Printf("user controller name:%s\n", c.Name)
}

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

func main() {
	members, err := digo.Members("group.controllers")

	if err == nil {
		for _, member := range members {
			member.(Controller).Print()
		}
	}
}
