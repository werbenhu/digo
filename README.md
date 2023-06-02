
<div align='center'>
<a href="https://github.com/werbenhu/digo/actions"><img src="https://github.com/werbenhu/digo/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/digo"><img src="https://goreportcard.com/badge/github.com/werbenhu/digo"></a>
<a href="https://coveralls.io/github/werbenhu/digo?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/digo/badge.svg?branch=master"></a>  
<a href="https://github.com/werbenhu/digo"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/digo"><img src="https://pkg.go.dev/badge/github.com/werbenhu/digo.svg"></a>
</div>

[English](README.md) | [简体中文](README-CN.md)

# digo

> Note: Go version `1.20+` is required.


**An annotation-based tool for compile-time dependency injection in Golang.**

## Features
- Use annotations in comments.
- Automatic code generation.
- Automatic detection of circular dependencies.
- Compile-time dependency injection.
- Automatic initialization.
- Support for managing instance groups.

## Quick Start


For more examples, please refer to: [examples](examples).

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
```

### Install digogen Tool

```sh
go install github.com/werbenhu/digo/digogen@v1.0.2
```

### Generate Dependency Injection Code

Open the command line and execute the following command. digogen will automatically generate the `digo.generated.go` source code file based on the annotations.
```sh
cd examples/simple
digogen
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