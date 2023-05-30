package database

import "fmt"

// @provider({"id": "database.mysql.url"})
func NewMysqlUrl() string {
	return "mysql:192.168.1.1:3306"
}

type Database interface {
	Print()
}

type Mysql struct {
	Url string
}

// @provider({"id": "database.mysql"})
// @inject({"param":"url", "id":"database.mysql.url"})
func NewMysql(url string) *Mysql {
	return &Mysql{
		Url: url,
	}
}

func (m *Mysql) Print() {
	fmt.Printf("Mysql Print url:%s\n", m.Url)
}
