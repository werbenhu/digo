package main

import (
	"github.com/werbenhu/digo"
	"github.com/werbenhu/digo/examples/multipackage/controllers"
	"github.com/werbenhu/digo/examples/multipackage/models"
)

func main() {
	user, err := digo.Provide("model.user")
	if err == nil {
		user.(*models.User).Print()
	}

	members, err := digo.Members("group.controllers")

	if err == nil {
		for _, member := range members {
			member.(controllers.Controller).Print()
		}
	}
}
