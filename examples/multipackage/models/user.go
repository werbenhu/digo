package models

import (
	"fmt"

	"github.com/werbenhu/digo/examples/multipackage/database"
)

type User struct {
	Db database.Database
}

// @provider({"id": "model.user"})
// @inject({"param":"db", "id":"database.mysql"})
func NewUser(db database.Database) *User {
	return &User{
		Db: db,
	}
}

func (u *User) Print() {
	fmt.Printf("User Print\n")
	u.Db.Print()
}
