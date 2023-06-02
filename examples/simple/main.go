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

type Redis struct {
}

// @provider({"id":"main.redis"})
func NewRedis() *Redis {
	return &Redis{}
}

type App struct {
	db    *Db
	redis *Redis
}

// @provider({"id":"main.app"})
// @inject({"param":"db", "id":"main.db"})
// @inject({"param":"redis", "id":"main.redis"})
func NewApp(db *Db, redis *Redis) *App {
	return &App{
		db:    db,
		redis: redis,
	}
}

func (a *App) Start() {
	log.Printf("app strat, db:%s\n", a.db.url)
}

func main() {
	app, err := digo.Provide("main.app")
	if err == nil {
		app.(*App).Start()
	}
}
