package tools

import (
	"fmt"

	"github.com/mochi-co/mqtt/v2"
	"github.com/werbenhu/eventbus"
	eb "github.com/werbenhu/eventbus"
)

type Tools struct {
	Bus *eb.EventBus
}

// @provider({"id":"name"})
// @group({"id":"strings"})
func NewName() string {
	fmt.Println("aaa")
	return "werben"
}

// @provider({"id":"class"})
// @group({"id":"strings"})
func NewClass() string {
	return "class"
}

// @provider({"id":"tools.eventbus"})
func NewEventBus() *eb.EventBus {
	return eb.New()
}

// @provider({"id":"tools.mq"})
func NewMqtt() *mqtt.Server {
	return nil
}

// @provider({"id":"tools.age"})
func NewAge() int {
	return 28
}

// @provider({"id":"tools.tools"})
// @inject({"param":"name", "id":"name"})
// @inject({"param":"class", "id":"class"})
// @inject({"param":"age", "id":"tools.age"})
// @inject({"param":"bus", "id":"tools.eventbus"})
// @inject({"param":"mq", "id":"tools.mq", "pkg": "github.com/mochi-co/mqtt/v2"})
// @group({"id":"tools"})
func NewTools(name, class string, age int, bus *eb.EventBus, mq *mqtt.Server) *Tools {
	return &Tools{}
}

// @provider({"id":"tools.tools2"})
// @inject({"param":"name", "id":"name"})
// @inject({"param":"class", "id":"class"})
// @inject({"param":"bus", "id":"tools.eventbus"})
func NewTools2(name, class string, bus eventbus.EventBus, md5 *Md5) *Tools {
	return &Tools{}
}
