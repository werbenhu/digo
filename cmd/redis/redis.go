package redis

import (
	t "github.com/werbenhu/digo/cmd/tools"
)

type Redis struct {
}

// @provider({"id":"redis"})
// @inject({"param":"name", "id":"name"})
// @inject({"param":"tools", "id":"tools.tools"})
func NewRedis(name string, tools t.Tools) *Redis {
	return &Redis{}
}
