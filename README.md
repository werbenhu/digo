# A Golang Dependency Injection Tool


[English](README.md) | [简体中文](README-CN.md)

## Quick Start

### Write Code and Annotations

```go
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
	log.Printf("app start, db:%s\n", a.Db.url)
}

func main() {
	app, err := digo.Provide("main.app")
	if err == nil {
		app.(*App).Start()
	}
}
```

### Install digocli Tool

```sh
go install github.com/werbenhu/digo/digocli@v0.0.1
```

### Generate Dependency Injection Code

Open the command line and execute the following command. digocli will automatically generate the `digo.generated.go` source code file based on the annotations.
```sh
digocli
```

### Run the Code

`go run .\digo.generated.go .\main.go`


## Annotation Details

### @provider

The `@provider` annotation indicates that it is an instance provider, and the instance is a singleton.

- Example:
```
// @provider({"id":"main.db"})
```
- Supported parameters:

| Name | Type | Required | Description |
| -------- | -----: | -----: | :----: |
| id     | string |  Yes| The ID of the instance    |

To obtain an instance, you can use digo.Provide(providerId) to retrieve the instance of a specific provider.
```go
app, err := digo.Provide("main.app")
if err == nil {
	app.(*App).Start()
}

```

## @inject
The `@inject` annotation indicates injecting an instance into a parameter. The `@inject` annotation must coexist with either `@provider` or `@group`.

- Example:
```
// @inject({"param":"db", "id":"main.db"})
```
- Supported parameters:

| Name | Type | Required | Description   |
| -------- | -----:  | -----:  |:----:  |
| param     | string |Yes|   Specifies the parameter to inject the instance into    |
| id     | string | Yes|   Specifies the ID of the instance to be injected    |
| pkg     | string | No |   Specifies the package to import for the parameter    |

The `pkg` parameter is used when you need to import a specific package. For example, if you need to import the package `github.com/xxx/tool/v1`, you would use the package name as `*tool.Struct`, not `*v1.Struct`. In such cases, you need to explicitly specify the import of the `github.com/xxx/tool/v1` package.

```
// @inject({"param":"tool", "id":"main.tool", "pkg":"github.com/xxx/tool/v1"})
```

## @group

The `@group` annotation indicates registering an instance to a group.

- Example:
```
// @group({"id":"main.controllers"})
```

- Supported parameters:

| Name | Type | Required | Description |
| -------- | -----: | -----: | :----: |
| id     | string |  Yes | The ID of the group   |

To retrieve all instances of a group, you can use `digo.Members(groupId)` to get all the instances of the group.

```go
ctrls, err := digo.Members("main.controllers")
if err == nil {
    for _, controller := range ctrls {
        // TODO:
    }
}
```