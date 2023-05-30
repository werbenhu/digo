package controllers

import "fmt"

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
