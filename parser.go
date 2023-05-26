package digo

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	// RegexpText represents the regular expression pattern for parsing annotations.
	RegexpText = `^//\s*@(provider|inject|group)\s*\((.*)\s*\)`
)

// Provider represents a provider.
type Provider struct {
	Id string // Id represents the identifier of the provider.
}

// Member represents a member in a group.
type Member struct {
	GroupId string `json:"id"` // GroupId represents the group ID of the member.
}

// Injector represents an injector parameter.
type Injector struct {
	ProviderId string `json:"id"` // Id represents the identifier of the provider.
	Pkg        string
	Param      string // Param represents the parameter name.
	Alias      string
	typ        ast.Expr // Typ represents the type of the parameter.
}

// GetObjName returns the object name of the injector.
func (i *Injector) GetObjName() string {
	return i.Param + "_obj"
}

// GetArgName returns the argument name of the injector.
func (i *Injector) GetArgName() string {
	return i.Param
}

type DiFunc struct {
	name      string
	injectors []*Injector
	provider  string
	group     string
	sort      int
	pkg       *DiPackage
	file      *DiFile
}

func replaceSeparator(id string) string {
	name := strings.ReplaceAll(id, ".", "_")
	name = strings.ReplaceAll(name, "/", "_")
	return name
}

func NewDiFunc(pkg *DiPackage, file *DiFile, name string) *DiFunc {
	return &DiFunc{
		name:      name,
		pkg:       pkg,
		file:      file,
		injectors: make([]*Injector, 0),
	}
}

func (fn *DiFunc) providerArgName() string {
	return replaceSeparator(fn.provider)
}

func (fn *DiFunc) providerObjName() string {
	return replaceSeparator(fn.provider) + "_obj"
}

func (fn *DiFunc) providerFuncName() string {
	return "init_" + replaceSeparator(fn.provider)
}

func (fn *DiFunc) groupFuncName() string {
	return "group_" + replaceSeparator(fn.group) + "_" + fn.name
}

type DiImport struct {
	name string // Alias represents the package alias.
	path string // Path represents the import path.
}

type DiFile struct {
	name    string
	pkg     *DiPackage
	imports map[string]*DiImport
}

func NewDifile(pkg *DiPackage, name string) *DiFile {
	return &DiFile{
		name:    name,
		pkg:     pkg,
		imports: make(map[string]*DiImport),
	}
}

type DiPackage struct {
	name   string
	path   string
	folder string

	funcs map[string]*DiFunc
	files map[string]*DiFile
}

func NewDiPackage(name string, path string, folder string) *DiPackage {
	return &DiPackage{
		name:   name,
		path:   path,
		folder: folder,
		funcs:  make(map[string]*DiFunc),
		files:  make(map[string]*DiFile),
	}
}

func (pkg *DiPackage) findProvider(id string) *DiFunc {
	for _, fn := range pkg.funcs {
		if id == fn.provider {
			return fn
		}
	}
	return nil
}

type Parser struct {
	packages []*DiPackage
}

func NewParser() *Parser {
	return &Parser{
		packages: make([]*DiPackage, 0),
	}
}

// parseImports analyzes and extracts information about imported packages in the current file.
func (p *Parser) parseImports(pkg *DiPackage, file *DiFile, decl *ast.GenDecl) {
	if decl.Tok == token.IMPORT {
		for _, spec := range decl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				return
			}
			if importSpec.Name != nil && len(importSpec.Name.Name) > 0 {
				file.imports[importSpec.Name.Name] = &DiImport{
					name: importSpec.Name.Name,
					path: strings.ReplaceAll(importSpec.Path.Value, "\"", ""),
				}
			} else {
				splitted := strings.Split(importSpec.Path.Value, "/")
				name := strings.ReplaceAll(splitted[len(splitted)-1], "\"", "")
				file.imports[name] = &DiImport{
					path: strings.ReplaceAll(importSpec.Path.Value, "\"", ""),
				}
			}
		}
	}
}

// matchComment 匹配符合provider， inject， group规则的注释
// 返回注解类型以及注解json格式的内容
// matchComment matches comments that comply with the provider, inject, group rules.
// It returns the annotation type and the JSON-formatted content of the annotation.
func (p *Parser) matchComment(comment string) (name string, body string) {
	r := regexp.MustCompile(RegexpText)
	if matches := r.FindStringSubmatch(comment); matches != nil {
		name = matches[1]
		body = matches[2]
	}
	return
}

// findProvider 根据id查找provider
// findProviderfinds a provider by its ID.
func (p *Parser) findProvider(id string) *DiFunc {
	for _, pkg := range p.packages {
		fn := pkg.findProvider(id)
		if fn != nil {
			return fn
		}
	}
	return nil
}

// parseProvider 分析提取源码中所有的@provider注解，并将注解信息保存在Provider对象中
// parseProvider analyzes and extracts all the @provider annotations in the source code and saves the annotation information in a Provider object.
func (p *Parser) parseProvider(body string, fn *DiFunc) error {
	provider := &Provider{}
	if err := json.Unmarshal([]byte(body), provider); err != nil {
		return fmt.Errorf("wrong json format, %s", err.Error())
	}

	if p.findProvider(provider.Id) != nil || fn.pkg.findProvider(provider.Id) != nil {
		return fmt.Errorf("[ERROR] duplicate defined provider id: %s", provider.Id)
	}
	fn.provider = provider.Id
	return nil
}

// parseInject 分析源码码中所有的@inject注解，并将inject信息提取到Injector对象中
// parseInject analyzes all the @inject annotations in the source code and extracts the inject information into an Injector object.
func (p *Parser) parseInject(body string, fn *DiFunc, decl *ast.FuncDecl) error {
	injector := &Injector{}
	if err := json.Unmarshal([]byte(body), injector); err != nil {
		return err
	}

	// @inject注解里的param必须在函数的参数列表中要能找到
	// isFieldFound标识是否在函数参数列表中存在该参数
	// The "param" in the @inject annotation must be found in the function's parameter list.
	// isFieldFound indicates whether the parameter exists in the function's parameter list.
	isFieldFound := false

	for i, field := range decl.Type.Params.List {
		for _, name := range field.Names {
			if name.Name == injector.Param {
				isFieldFound = true
				injector.typ = decl.Type.Params.List[i].Type
			}
		}
	}

	// 如果没有找到注解中的param在函数中对应的参数，则返回一个错误
	// If the parameter specified in the annotation is not found in the function's parameters, return an error.
	if !isFieldFound {
		return errors.New("injected parameter is not found")
	}

	// @inject注解可以显式的设置变量对应的包名，
	// 比如@inject({"param":"mq", "id":"mq", "pkg": "github.com/mochi-co/mqtt/v2"})
	// 这里的pkg就是这个变量需要引入的包
	// 如果@inject注解里显式表明了该param需要引入的包，则不需要再去当前文件的引入包列表中去查找了

	// The @inject annotation can explicitly specify the package name for the variable,
	// e.g., @inject({"param":"mq", "id":"mq", "pkg": "github.com/mochi-co/mqtt/v2"})
	// In this case, the package does not need to be searched in the current file's import list.
	if len(injector.Pkg) > 0 {
		fn.injectors = append(fn.injectors, injector)
		return nil
	}

	// 注入的参数，如果该参数的类型不是当前包下定义的，需要引入别的包
	// 比如需要注入一个参数类型是: eventbus.EventBus, 这个类型是包github.com/werbenhu/eventbus里定义的
	// 这里需要从当前文件的import列表中，找出这个包名

	// For injected parameters whose types are not defined in the current package, this packages need to be imported.
	// For example, if a parameter with type eventbus.EventBus is injected,
	// and this type is defined in the package github.com/werbenhu/eventbus,
	// the package name needs to be found from the import list of the current file.
	var importPkg string

	// pkgFound标识是否从当前文件的import列表中找到了该变量需要使用的包
	// isPkgFound indicates whether the package required by the variable has been found in the current go file's import list.
	var isPkgFound bool

	if starExpr, ok := injector.typ.(*ast.StarExpr); ok {
		if selExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {

			// 如果是指针类型参数，类似 *aaa.bbb这种类型
			// If the parameter type is a pointer type, e.g., *aaa.bbb
			importPkg = selExpr.X.(*ast.Ident).Name
			if selExpr.X == nil {

				// 如果是 *struct这种类型的参数，说明这个参数使用的当前包中定义的struct，
				// 这种类型不需要额外的引入包了
				// If it is a *struct type parameter, it means that the parameter uses a struct defined in the current package,
				// and no additional package needs to be imported.
				isPkgFound = true

			} else {
				// 如果是*pkg.struct这种类型的参数，说明这个参数需要引入包
				// 从当前文件的import列表中，将需要引入的包找出来，并放在injector中
				// 如果引入包的时候使用了别名，那么这里别名也需要保存
				// If it is a *pkg.struct type parameter, it means that the parameter requires importing a package.
				// Find the required package from the import list of the current file and store it in the injector.
				// If an alias is used for the imported package, the alias needs to be stored as well.
				if impor, ok := fn.file.imports[importPkg]; ok {
					injector.Pkg = impor.path
					injector.Alias = impor.name
					isPkgFound = true
				}
			}
		} else if _, ok := starExpr.X.(*ast.Ident); ok {

			// 如果是指针类型参数，类似 *aaa这种类型
			// 说明该类型是在当前包定义的，不需要引用其他包
			// If it is a pointer type parameter, e.g., *aaa
			// It means that the type is defined in the current package, and no other package needs to be imported.
			isPkgFound = true
		}
	} else if selExpr, ok := injector.typ.(*ast.SelectorExpr); ok {
		importPkg = selExpr.X.(*ast.Ident).Name

		// 如果是pkg.struct这种类型的参数，说明这个参数需要引入特殊的包
		// 从当前文件的import列表中，将需要引入的包找出来，并放在injector中
		// 如果引入包的时候使用了别名，那么这里别名也需要保存

		// If it is a pkg.struct type parameter, it means that the parameter requires importing a specific package.
		// Find the required package from the import list of the current file and store it in the injector.
		// If an alias is used for the imported package, the alias needs to be stored as well.
		if selExpr.X != nil {
			if impor, ok := fn.file.imports[importPkg]; ok {
				injector.Pkg = impor.path
				injector.Alias = impor.name
				isPkgFound = true
			}
		}
	} else if _, ok := injector.typ.(*ast.Ident); ok {
		isPkgFound = true
	}

	if !isPkgFound {
		return errors.New("injected parameter's package not found")
	}
	fn.injectors = append(fn.injectors, injector)
	return nil
}

// parseGroup 分析提取源码中所有的@group注解，并将注解信息保存在Provider对象中
// parseGroup analyzes and extracts all the @group annotations in the source code and saves the annotation information in a Member object.
func (p *Parser) parseGroup(body string, fn *DiFunc) error {
	member := &Member{}
	if err := json.Unmarshal([]byte(body), member); err != nil {
		return err
	}
	fn.group = member.GroupId
	return nil
}

// parseFunc 分析某一个函数的注解，提取出provider、inject、group信息
// parseFunc analyzes the annotations of a specific function and extracts provider, inject, and group information.
func (p *Parser) parseFunc(pkg *DiPackage, fn *DiFunc, decl *ast.FuncDecl) error {

	// 如果注释不为空
	// If the comment is not empty
	if decl.Doc != nil && decl.Doc.List != nil {
		for _, comment := range decl.Doc.List {
			// 正则表达式匹配注解规则
			// Use regular expressions to match annotation rules
			name, body := p.matchComment(comment.Text)
			switch name {
			case "provider":
				if err := p.parseProvider(body, fn); err != nil {
					return fmt.Errorf("failed to parse provider annotation, %s in package: %s Func: %s", err.Error(), pkg.path, fn.name)
				}
			case "inject":
				if err := p.parseInject(body, fn, decl); err != nil {
					return fmt.Errorf("failed to parse inject annotation, %s in package: %s Func: %s", err.Error(), pkg.path, fn.name)
				}
			case "group":
				if err := p.parseGroup(body, fn); err != nil {
					return fmt.Errorf("failed to parse group annotation, %s in package: %s Func: %s", err.Error(), pkg.path, fn.name)
				}
			}
		}
	}

	// 检查是否还有参数没有被注入
	// Check if there are any parameters that have not been injected.
	for _, field := range decl.Type.Params.List {
		for _, name := range field.Names {
			found := false
			for _, injector := range fn.injectors {
				if name.String() == injector.Param {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("all parameters of the provider must be injected, params: %v have not been injected yet, in package: %s Func: %s\n",
					name.String(), pkg.path, fn.name)
			}
		}
	}
	return nil
}

// parse 解析所有包下函数的注释，并且提取出每个文件的import的包的信息
// parse parses the comments of functions in all packages and extracts the import package information for each file.
func (p *Parser) parse(pkgs []*packages.Package) error {
	for _, pkg := range pkgs {

		splitted := strings.Split(pkg.GoFiles[0], string(os.PathSeparator))
		folder := strings.Join(splitted[:len(splitted)-1], string(os.PathSeparator))
		diPkg := NewDiPackage(pkg.Name, pkg.PkgPath, folder)

		for _, syntax := range pkg.Syntax {
			diFile := NewDifile(diPkg, syntax.Name.String())
			for _, decl := range syntax.Decls {

				if genDecl, ok := decl.(*ast.GenDecl); ok {
					p.parseImports(diPkg, diFile, genDecl)
				} else if fn, ok := decl.(*ast.FuncDecl); ok {

					diFunc := NewDiFunc(diPkg, diFile, fn.Name.String())
					if err := p.parseFunc(diPkg, diFunc, fn); err != nil {
						return err
					}

					if len(diFunc.provider) > 0 || len(diFunc.group) > 0 {
						diPkg.funcs[diFunc.name] = diFunc
					}
				}
			}
			diPkg.files[syntax.Name.String()] = diFile
		}

		if len(diPkg.funcs) > 0 {
			p.packages = append(p.packages, diPkg)
		}
	}
	return nil
}

// Start 启动分析注解，并生成go代码，写入到文件中
// Start initiates the analysis of annotations, generates Go code, and writes it to files.
func (p *Parser) Start() {
	// Load packages and their syntax.
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
	}, "pattern=./...")

	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// Parse annotations and extract information.
	if err := p.parse(pkgs); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return
	}

	// Check the legality of injectors and cyclic provider dependencies.
	// if p.checkInjectorLegal() && p.checkCyclicProvider() {
	// Generate Go code for initializing providers.
	for _, pkg := range p.packages {
		generator := NewGenerator(pkg)
		generator.Do()
	}
	// }
}
