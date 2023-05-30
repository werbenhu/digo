package main

import (
	"log"

	"github.com/werbenhu/digo"
)

// @provider({"id":"main.db.url"})
func NewDbUrl() string {
	return "localhost:3306"
}

type Db struct {
	url string
}

// @provider({"id":"main.db"})
// @inject({"param":"url", "id":"main.db.url"})
func NewDb(url string) *Db {
	return &Db{
		url: url,
	}
}

type App struct {
	Db *Db
}

// @provider({"id":"main.app"})
// @inject({"param":"db", "id":"main.db"})
func NewApp(db *Db) *App {
	return &App{
		Db: db,
	}
}

func (a *App) Start() {
	log.Printf("app strat, db:%s\n", a.Db.url)
}

func main() {
	app, err := digo.Provide("main.app")
	if err == nil {
		app.(*App).Start()
	}
}
