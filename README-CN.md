

<div align='center'>
<a href="https://github.com/werbenhu/digo/actions"><img src="https://github.com/werbenhu/digo/workflows/Go/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/werbenhu/digo"><img src="https://goreportcard.com/badge/github.com/werbenhu/digo"></a>
<a href="https://coveralls.io/github/werbenhu/digo?branch=master"><img src="https://coveralls.io/repos/github/werbenhu/digo/badge.svg?branch=master"></a>  
<a href="https://github.com/werbenhu/digo"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
<a href="https://pkg.go.dev/github.com/werbenhu/digo"><img src="https://pkg.go.dev/badge/github.com/werbenhu/digo.svg"></a>
</div>

[English](README.md) | [简体中文](README-CN.md)

# digo
> 注意: digo需要Go 版本 `1.20+`.

**一个通过注解来实现依赖注入的golang工具**

## 特性

- 使用注释中的注解
- 自动检测循环依赖
- 编译时期依赖注入
- 自动初始化
- 支持实例组的管理


## 快速开始

更多示例请参考：[examples](examples)

### 编写代码和注解

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
	log.Printf("app strat, db:%s\n", a.Db.url)
}

func main() {
	app, err := digo.Provide("main.app")
	if err == nil {
		app.(*App).Start()
	}
}
```

### 安装digogen工具

```sh
go install github.com/werbenhu/digo/digogen@v1.0.2
```
### 生成依赖注入代码
打开命令行执行下面命令,`digogen`将会根据注解自动生成`digo.generated.go`源码文件.
```sh
digogen
```
### 运行代码
`go run  .\digo.generated.go .\main.go`

## 注解详情

### @provider
@provider注解表示是一个实例提供者，该实例是一个单例
- 示例
```
// @provider({"id":"main.db"})
```
- 支持的参数：

| 参数 | 类型 | 是否必需 | 说明 |
| -------- | -----: | -----: | :----: |
| id     | string |  是| 实例的id    |

如果获取实例，通过`digo.Provide(providerId)`可以获取到某一个provider的实例
```
app, err := digo.Provide("main.app")
if err == nil {
	app.(*App).Start()
}
```

### @inject
@inject注解表示注入一个实例到某个参数, @inject注解必须和@provider或者@group二者中的一个同时存在.
- 示例
```
// @inject({"param":"db", "id":"main.db"})
```
- 支持的参数：

| 参数 | 类型 | 是否必需 | 说明  |
| -------- | -----:  | -----:  |:----:  |
| param     | string |是|   指明哪个参数需要注入实例    |
| id     | string | 是|   指明需要注入的实例id    |
| pkg     | string | 否 |   该参数需要引入特定的包    |

pkg在什么时候需要使用，比如我们需要引入一个包 `github.com/xxx/tool/v1` , 我们使用包名的时候是这样使用的 *tool.Struct， 而不是 *v1.Struct，那我们需要显示指明需要导入`github.com/xxx/tool/v1`包

```
// @inject({"param":"tool", "id":"main.tool", "pkg":"github.com/xxx/tool/v1"})
```

### @group
@group注解表示将实例注册到一个组
- 示例
```
// @group({"id":"main.controllers"})
```
- 支持的参数：

| 参数 | 类型 | 是否必需 | 说明 |
| -------- | -----: | -----: | :----: |
| id     | string |  是| 组的id    |

如果获取组的所有实例，通过`digo.Members(groupId)`可以获取到组的所有实例
```
ctrls, err := digo.Members("main.controllers")
if err == nil {
    for _, controller := range ctrls {
        // TODO:
    }
}
```