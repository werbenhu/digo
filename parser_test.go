// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package digo

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

var (
	emptyType = &ast.FuncType{
		Params: &ast.FieldList{
			List: make([]*ast.Field, 0),
		},
	}

	declRegular = &ast.FuncDecl{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{
						newIdent("myParam"),
					},
					Type: newIdent("string"),
				}},
			},
		},
	}
	declCompound = &ast.FuncDecl{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{
						newIdent("myParam"),
					},
					Type: newSelectorExpr("pkg.string"),
				}},
			},
		},
	}
	declCompoundPointer = &ast.FuncDecl{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{
						newIdent("myParam"),
					},
					Type: newStarExpr("pkg.string"),
				}},
			},
		},
	}
)

func TestChain_Insert(t *testing.T) {
	chain := newChain()
	fn1 := &DiFunc{ProviderId: "provider1"}
	fn2 := &DiFunc{ProviderId: "provider2"}
	fn3 := &DiFunc{ProviderId: "provider3"}

	// Inserting the first provider should succeed.
	if !chain.insert(fn1) {
		t.Error("Expected insert to succeed for the first provider")
	}

	// Inserting the second provider should succeed.
	if !chain.insert(fn2) {
		t.Error("Expected insert to succeed for the second provider")
	}

	// Inserting the first provider again should fail, indicating a cyclic dependency.
	if chain.insert(fn1) {
		t.Error("Expected insert to fail for the first provider indicating cyclic dependency")
	}

	// Inserting the third provider should succeed.
	if !chain.insert(fn3) {
		t.Error("Expected insert to succeed for the third provider")
	}
}

func TestChain_String(t *testing.T) {
	chain := chain{
		&DiFunc{ProviderId: "provider1"},
		&DiFunc{ProviderId: "provider2"},
		&DiFunc{ProviderId: "provider3"},
	}

	expected := "provider1 -> provider2 -> provider3"
	result := chain.String()
	require.Equal(t, result, expected)
}

func TestInjector_GetObjName(t *testing.T) {
	injector := &Injector{
		Param: "param",
	}

	expected := "param_obj"
	result := injector.GetObjName()
	require.Equal(t, result, expected)
}

func TestInjector_GetArgName(t *testing.T) {
	injector := &Injector{
		Param: "param",
	}

	expected := "param"
	result := injector.GetArgName()
	require.Equal(t, result, expected)
}

func TestReplaceSeparator(t *testing.T) {
	id := "provider.id/example/path"

	expected := "provider_id_example_path"
	result := replaceSeparator(id)
	require.Equal(t, result, expected)

	id = "provider.id/another.path"

	expected = "provider_id_another_path"
	result = replaceSeparator(id)
	require.Equal(t, result, expected)
}

func TestDiFunc_providerArgName(t *testing.T) {
	fn := &DiFunc{
		ProviderId: "provider.id",
	}

	expected := "provider_id"
	result := fn.providerArgName()
	require.Equal(t, result, expected)
}

func TestDiFunc_providerObjName(t *testing.T) {
	fn := &DiFunc{
		ProviderId: "provider.id",
	}

	expected := "provider_id_obj"
	result := fn.providerObjName()
	require.Equal(t, result, expected)
}

func TestDiFunc_providerFuncName(t *testing.T) {
	fn := &DiFunc{
		ProviderId: "provider.id",
	}

	expected := "init_provider_id"
	result := fn.providerFuncName()
	require.Equal(t, result, expected)
}

func TestDiFunc_groupFuncName(t *testing.T) {
	fn := &DiFunc{
		GroupId: "group.id",
		Name:    "example",
	}

	expected := "group_group_id_example"
	result := fn.groupFuncName()
	require.Equal(t, result, expected)
}

func TestDiFuncs_Len(t *testing.T) {
	funcs := DiFuncs{
		&DiFunc{},
		&DiFunc{},
		&DiFunc{},
	}

	expected := 3
	result := funcs.Len()

	assert.Equal(t, expected, result, "Unexpected length")
}

func TestDiFuncs_Swap(t *testing.T) {
	funcs := DiFuncs{
		&DiFunc{Name: "A"},
		&DiFunc{Name: "B"},
	}

	funcs.Swap(0, 1)

	expected := DiFuncs{
		&DiFunc{Name: "B"},
		&DiFunc{Name: "A"},
	}

	assert.Equal(t, expected, funcs, "Elements not swapped")
}

func TestDiFuncs_Less(t *testing.T) {
	funcs := DiFuncs{
		&DiFunc{Sort: 3},
		&DiFunc{Sort: 2},
	}

	result := funcs.Less(0, 1)
	assert.True(t, result, "Expected element at position 0 to be less than element at position 1")
}

func TestDiFuncs_Sort(t *testing.T) {
	funcs := DiFuncs{
		&DiFunc{Sort: 2},
		&DiFunc{Sort: 1},
		&DiFunc{Sort: 3},
	}
	funcs.Sort()

	expected := DiFuncs{
		&DiFunc{Sort: 3},
		&DiFunc{Sort: 2},
		&DiFunc{Sort: 1},
	}
	assert.Equal(t, expected, funcs, "Elements not sorted in reverse order")
}

func TestNewDiFile(t *testing.T) {
	pkg := &DiPackage{}
	name := "example.go"
	diFile := NewDiFile(pkg, name)

	assert.Equal(t, name, diFile.Name, "Unexpected file name")
	assert.Equal(t, pkg, diFile.Package, "Unexpected package")
	assert.NotNil(t, diFile.Imports, "Imports map should not be nil")
}

func TestNewDiPackage(t *testing.T) {
	name := "example"
	path := "/path/to/example"
	folder := "/path/to/folder"
	diPackage := NewDiPackage(name, path, folder)

	assert.Equal(t, name, diPackage.Name, "Unexpected package name")
	assert.Equal(t, path, diPackage.Path, "Unexpected package path")
	assert.Equal(t, folder, diPackage.Folder, "Unexpected package folder")
	assert.NotNil(t, diPackage.Funcs, "Funcs slice should not be nil")
	assert.NotNil(t, diPackage.Files, "Files map should not be nil")
}

func TestDiPackage_FindProvider(t *testing.T) {
	name := "example"
	path := "/path/to/example"
	folder := "/path/to/folder"
	diPackage := NewDiPackage(name, path, folder)

	providerID := "provider1"
	diFunc := &DiFunc{ProviderId: providerID}
	diPackage.Funcs = append(diPackage.Funcs, diFunc)

	result := diPackage.findProvider(providerID)
	assert.Equal(t, diFunc, result, "Unexpected provider found")
}

func TestDiPackage_FindProvider_NotFound(t *testing.T) {
	name := "example"
	path := "/path/to/example"
	folder := "/path/to/folder"
	diPackage := NewDiPackage(name, path, folder)

	providerID := "provider1"
	result := diPackage.findProvider(providerID)
	assert.Nil(t, result, "Expected nil provider")
}

func TestNewParser(t *testing.T) {
	parser := NewParser()

	assert.NotNil(t, parser.Packages, "Packages slice should not be nil")
	assert.Empty(t, parser.Packages, "Packages slice should be empty")
}

func TestParser_parseImports(t *testing.T) {
	pkg := NewDiPackage("example", "/path/to/example", "/path/to/folder")
	file := NewDiFile(pkg, "example.go")
	decl := &ast.GenDecl{Tok: token.IMPORT}
	parser := &Parser{}

	parser.parseImports(pkg, file, decl)
	assert.NotNil(t, file.Imports, "Imports map should not be nil")
	assert.Empty(t, file.Imports, "Imports map should be empty")
}

func TestParser_parseImports_WithNamedImport(t *testing.T) {
	pkg := NewDiPackage("example", "/path/to/example", "/path/to/folder")
	file := NewDiFile(pkg, "example.go")
	decl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Name: &ast.Ident{Name: "alias"},
				Path: &ast.BasicLit{Value: "\"example.com/package\""},
			},
		},
	}
	parser := &Parser{}

	parser.parseImports(pkg, file, decl)
	expectedImport := &DiImport{Name: "alias", Path: "example.com/package"}
	assert.Equal(t, expectedImport, file.Imports["alias"], "Unexpected named import")
}

func TestParser_parseImports_WithoutNamedImport(t *testing.T) {
	pkg := NewDiPackage("example", "/path/to/example", "/path/to/folder")
	file := NewDiFile(pkg, "example.go")
	decl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Name: nil,
				Path: &ast.BasicLit{Value: "\"example.com/package\""},
			},
		},
	}

	parser := &Parser{}
	parser.parseImports(pkg, file, decl)
	splitted := strings.Split("\"example.com/package\"", "/")

	name := strings.ReplaceAll(splitted[len(splitted)-1], "\"", "")
	expectedImport := &DiImport{Path: "example.com/package"}
	assert.Equal(t, expectedImport, file.Imports[name], "Unexpected import without name")
}
func TestParser_MatchComment(t *testing.T) {
	parser := NewParser()

	// Test case 1: Matching comment
	comment := "// @provider({\"id\":\"provider1\"})"
	name, body := parser.matchComment(comment)

	expectedName := "provider"
	expectedBody := "{\"id\":\"provider1\"}"

	assert.Equal(t, expectedName, name, "Expected name to match")
	assert.Equal(t, expectedBody, body, "Expected body to match")

	// Test case 2: Non-matching comment
	comment = "// Some comment without annotation"
	name, body = parser.matchComment(comment)

	assert.Equal(t, "", name, "Expected name to be empty")
	assert.Equal(t, "", body, "Expected body to be empty")
}

func TestParser_FindProvider(t *testing.T) {
	parser := NewParser()
	pkg1 := NewDiPackage("pkg1", "path1", "folder1")
	pkg2 := NewDiPackage("pkg2", "path2", "folder2")
	fn1 := &DiFunc{
		ProviderId: "provider1",
	}
	fn2 := &DiFunc{
		ProviderId: "provider2",
	}
	pkg1.Funcs = []*DiFunc{fn1}
	pkg2.Funcs = []*DiFunc{fn2}
	parser.Packages = []*DiPackage{pkg1, pkg2}

	// Test case 1: Provider found
	foundProvider := parser.findProvider("provider2")
	assert.Equal(t, fn2, foundProvider, "Expected found provider to match")

	// Test case 2: Provider not found
	foundProvider = parser.findProvider("provider3")
	assert.Nil(t, foundProvider, "Expected provider not to be found")
}

func TestParser_ParseProvider(t *testing.T) {
	parser := NewParser()
	pkg := NewDiPackage("example", "/path/to/example", "/path/to/folder")
	file := NewDiFile(pkg, "example.go")
	fn := NewDiFunc(pkg, file, "ExampleFunc")

	// Test case 1: Valid JSON format
	body := "{\"id\":\"provider1\"}"
	err := parser.parseProvider(body, fn)
	expectedId := "provider1"
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expectedId, fn.ProviderId, "Expected provider ID to match")

	// Test case 2: Invalid JSON format
	body = "invalid json"
	err = parser.parseProvider(body, fn)
	expectedErrorMessage := "wrong JSON format: invalid character 'i' looking for beginning of value"
	assert.Error(t, err, "Expected error")
	assert.EqualError(t, err, expectedErrorMessage, "Expected error message to match")

	// Test case 3: Duplicate provider ID
	body = "{\"id\":\"provider1\"}"
	fn2 := NewDiFunc(pkg, file, "ExampleFunc1")
	fn2.ProviderId = "provider1"
	parser.Packages = append(parser.Packages, &DiPackage{Funcs: []*DiFunc{fn2}})
	err = parser.parseProvider(body, fn)

	expectedErrorMessage = "[ERROR] duplicate provider ID: provider1"
	assert.Error(t, err, "Expected error")
	assert.EqualError(t, err, expectedErrorMessage, "Expected error message to match")
}

func TestParser_ParseInject(t *testing.T) {
	parser := NewParser()

	// Test case 1: Valid JSON format, regular type
	body := "{\"param\":\"myParam\"}"
	fn := &DiFunc{
		File: &DiFile{
			Imports: map[string]*DiImport{
				"": {
					Name: "",
					Path: "",
				},
			},
		},
	}

	err := parser.parseInject(body, fn, declRegular)
	expectedParam := "myParam"
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 1, len(fn.Injectors), "Expected one injector")
	assert.Equal(t, expectedParam, fn.Injectors[0].Param, "Expected injector param to match")

	// Test case 2: Invalid JSON format
	body = "invalid json"
	err = parser.parseInject(body, fn, declRegular)
	expectedErrorMessage := "invalid character 'i' looking for beginning of value"
	assert.Error(t, err, "Expected error")
	assert.EqualError(t, err, expectedErrorMessage, "Expected error message to match")

	// Test case 3: Injected parameter not found
	body = "{\"param\":\"otherParam\"}"
	err = parser.parseInject(body, fn, declRegular)
	expectedErrorMessage = "injected parameter is not found"
	assert.Error(t, err, "Expected error")
	assert.EqualError(t, err, expectedErrorMessage, "Expected error message to match")

	// Test case 4: Package not found
	// body = "{\"param\":\"myParam\"}"
	// fn.File.Imports = map[string]*DiImport{}
	// err = parser.parseInject(body, fn, declCompoundPointer)
	// expectedErrorMessage = "injected parameter's package not found"
	// assert.Error(t, err, "Expected error")
	// assert.EqualError(t, err, expectedErrorMessage, "Expected error message to match")

	// Test case 5: Compound pointer type, package found
	body = "{\"param\":\"myParam\"}"
	fn.File.Imports["pkg"] = &DiImport{
		Name: "pkg",
		Path: "github.com/mochi-co/mqtt/v2",
	}
	err = parser.parseInject(body, fn, declCompoundPointer)

	expectedParam = "myParam"
	expectedPkg := "github.com/mochi-co/mqtt/v2"
	expectedAlias := "pkg"

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 2, len(fn.Injectors), "Expected two injectors")
	assert.Equal(t, expectedParam, fn.Injectors[1].Param, "Expected injector param to match")
	assert.Equal(t, expectedPkg, fn.Injectors[1].Pkg, "Expected injector package to match")
	assert.Equal(t, expectedAlias, fn.Injectors[1].Alias, "Expected injector alias to match")

	// Test case 6: Compound regular type, package not found
	body = "{\"param\":\"myParam\"}"
	fn.File.Imports = map[string]*DiImport{}
	err = parser.parseInject(body, fn, declCompound)
	expectedErrorMessage = "injected parameter's package not found"
	assert.Error(t, err, "Expected error")
	assert.EqualError(t, err, expectedErrorMessage, "Expected error message to match")
}

func TestParser_ParseGroup(t *testing.T) {
	parser := NewParser()
	pkg := NewDiPackage("example", "/path/to/example", "/path/to/folder")
	file := NewDiFile(pkg, "example.go")
	fn := NewDiFunc(pkg, file, "ExampleFunc")

	// Test case 1: Valid JSON format
	body := "{\"id\":\"myGroup\"}"
	err := parser.parseGroup(body, fn)
	expectedGroupId := "myGroup"
	assert.NoError(t, err)
	assert.Equal(t, expectedGroupId, fn.GroupId)

	// Test case 2: Invalid JSON format
	body = "invalid json"
	err = parser.parseGroup(body, fn)
	expectedErrorMessage := "invalid character 'i' looking for beginning of value"
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func TestParser_ParseFunc(t *testing.T) {

	parser := NewParser()
	pkg := NewDiPackage("example", "github.com/my/package", "/path/to/folder")
	file := NewDiFile(pkg, "example.go")
	fn := NewDiFunc(pkg, file, "ExampleFunc")

	comment := "// @provider({\"id\": \"myProvider\"})"
	funcType := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{{
				Names: []*ast.Ident{
					newIdent("myParam"),
				},
				Type: newIdent("string"),
			}},
		},
	}

	decl := &ast.FuncDecl{
		Doc:  newCommentGroup([]string{comment}),
		Type: emptyType,
	}

	// Test case 1: Provider annotation
	err := parser.parseFunc(pkg, fn, decl)
	expectedProviderId := "myProvider"
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expectedProviderId, fn.ProviderId)

	// Test case 2: Inject annotation
	injectComment := "// @inject({\"param\": \"myParam\"})"
	decl.Type = funcType
	decl.Doc = newCommentGroup([]string{comment, injectComment})
	fn = NewDiFunc(pkg, file, "myFunction")
	err = parser.parseFunc(pkg, fn, decl)

	expectedParam := "myParam"
	assert.NoError(t, err, "Expected no error")
	assert.Len(t, fn.Injectors, 1, "Expected one injector")
	assert.Equal(t, expectedParam, fn.Injectors[0].Param, "Expected injected parameter to match")

	// Test case 3: Group annotation
	comment = "// @group({\"id\": \"myGroup\"})"
	decl.Doc = newCommentGroup([]string{comment})
	decl.Type = emptyType
	fn = NewDiFunc(pkg, file, "myFunction")
	err = parser.parseFunc(pkg, fn, decl)

	expectedGroupId := "myGroup"
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expectedGroupId, fn.GroupId, "Expected group ID to match")

	// Test case 5: Missing injection
	comment = "// @provider({\"id\": \"myProvider\"})"
	decl.Doc = newCommentGroup([]string{comment})
	decl.Type = funcType

	pkg = NewDiPackage("example", "github.com/my/package", "/path/to/folder")
	fn = NewDiFunc(pkg, file, "myFunction")
	err = parser.parseFunc(pkg, fn, decl)

	expectedErrorMessage := "all parameters of the provider must be injected, param: myParam have not been injected yet, in pkg: github.com/my/package, function: myFunction"
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func TestParser_Parse(t *testing.T) {
	parser := NewParser()

	// Test case 1: Single package with functions and imports
	packagePath := "github.com/my/package"
	goFiles := []string{"file1.go", "file2.go"}

	pkg := &packages.Package{
		Name:    "package",
		PkgPath: packagePath,
		GoFiles: goFiles,
		Syntax: []*ast.File{{
			Name: newIdent("file1"),
			Decls: []ast.Decl{
				&ast.GenDecl{
					Tok: token.IMPORT,
					Specs: []ast.Spec{
						newImportSpec("github.com/other/package1", ""),
						newImportSpec("github.com/other/package2", ""),
					},
				},
			},
		}, {
			Name: newIdent("file2"),
			Decls: []ast.Decl{
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "myFunction",
					},
					Doc:  newCommentGroup([]string{"// @provider({\"id\": \"myProvider\"})"}),
					Type: emptyType,
				},
			},
		}},
	}

	err := parser.parse([]*packages.Package{pkg})
	assert.NoError(t, err, "Expected no error")
	assert.Len(t, parser.Packages, 1, "Expected one package")

	diPkg := parser.Packages[0]
	assert.Equal(t, packagePath, diPkg.Path, "Expected package path to match")
	assert.Equal(t, "package", diPkg.Name, "Expected package name to match")
	assert.Len(t, diPkg.Funcs, 1, "Expected one function")

	diFunc := diPkg.Funcs[0]
	assert.Equal(t, "myFunction", diFunc.Name, "Expected function name to match")
	assert.Equal(t, "myProvider", diFunc.ProviderId, "Expected provider ID to match")

	_, ok := diPkg.Files["file1"]
	assert.False(t, ok)
	_, ok = diPkg.Files["file2"]
	assert.True(t, ok)

	// Test case 2: Package without functions
	packagePath = "github.com/my/empty/package"
	goFiles = []string{"file.go"}
	pkg = &packages.Package{
		Name:    "empty",
		PkgPath: packagePath,
		GoFiles: goFiles,
		Syntax: []*ast.File{
			{
				Name: newIdent("file"),
				Decls: []ast.Decl{
					&ast.GenDecl{
						Tok: token.IMPORT,
						Specs: []ast.Spec{
							newImportSpec("github.com/other/package", ""),
						},
					},
				},
			},
		},
	}

	parser = NewParser()
	err = parser.parse([]*packages.Package{pkg})
	assert.NoError(t, err, "Expected no error")
	assert.Empty(t, parser.Packages, "Expected no packages")

	// Test case 3: Multiple packages
	package1Path := "github.com/my/package1"
	package2Path := "github.com/my/package2"
	package3Path := "github.com/my/package3"
	goFiles = []string{"file.go"}
	pkg1 := &packages.Package{
		Name:    "package1",
		PkgPath: package1Path,
		GoFiles: goFiles,
		Syntax: []*ast.File{
			{
				Name: newIdent("file"),
				Decls: []ast.Decl{
					&ast.GenDecl{
						Tok: token.IMPORT,
						Specs: []ast.Spec{
							newImportSpec("github.com/other/package", ""),
						},
					},
				},
			},
		},
	}
	pkg2 := &packages.Package{
		Name:    "package2",
		PkgPath: package2Path,
		GoFiles: goFiles,
		Syntax: []*ast.File{
			{
				Name: newIdent("file"),
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{
							Name: "myFunction",
						},
						Doc:  newCommentGroup([]string{"// @provider({\"id\": \"myProvider2\"})"}),
						Type: emptyType,
					},
				},
			},
		},
	}
	pkg3 := &packages.Package{
		Name:    "package3",
		PkgPath: package3Path,
		GoFiles: goFiles,
		Syntax: []*ast.File{
			{
				Name: newIdent("file"),
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{
							Name: "myFunction",
						},
						Doc:  newCommentGroup([]string{"// @provider({\"id\": \"myProvider3\"})"}),
						Type: emptyType,
					},
				},
			},
		},
	}

	err = parser.parse([]*packages.Package{pkg1, pkg2, pkg3})
	assert.NoError(t, err, "Expected no error")
	assert.Len(t, parser.Packages, 2, "Expected two packages")

	diPkg2 := parser.Packages[0]
	diPkg3 := parser.Packages[1]

	assert.Equal(t, package2Path, diPkg2.Path, "Expected package2 path to match")
	assert.Equal(t, "package2", diPkg2.Name, "Expected package2 name to match")
	assert.Len(t, diPkg2.Funcs, 1, "Expected one function in package2")

	assert.Equal(t, package3Path, diPkg3.Path, "Expected package3 path to match")
	assert.Equal(t, "package3", diPkg3.Name, "Expected package3 name to match")
	assert.Len(t, diPkg3.Funcs, 1, "Expected one function in package3")

	diFunc2 := diPkg2.Funcs[0]

	assert.Equal(t, "myFunction", diFunc2.Name, "Expected function name to match")
	assert.Equal(t, "myProvider2", diFunc2.ProviderId, "Expected provider ID to match")
}

func TestParser_FindProviderById(t *testing.T) {
	// Test case 1: Provider found
	parser := &Parser{
		Packages: []*DiPackage{
			{
				Funcs: []*DiFunc{{
					ProviderId: "myProvider",
				}},
			},
		},
	}

	result := parser.findProviderById("myProvider")
	assert.NotNil(t, result, "Expected provider to be found")
	assert.Equal(t, "myProvider", result.ProviderId, "Expected provider ID to match")

	// Test case 2: Provider not found
	result = parser.findProviderById("nonExistentProvider")
	assert.Nil(t, result, "Expected provider to not be found")
}

func TestParser_CheckInjectorLegal(t *testing.T) {
	// Test case 1: All injectors are legal
	func1 := &DiFunc{
		ProviderId: "provider1",
		Injectors: []*Injector{{
			ProviderId: "provider1",
		}},
	}
	func2 := &DiFunc{
		ProviderId: "provider2",
		Injectors: []*Injector{{
			ProviderId: "provider2",
			Dependency: func1,
		}},
	}

	parser := &Parser{
		Packages: []*DiPackage{{
			Funcs: []*DiFunc{func1, func2},
		}},
	}

	result := parser.checkInjectorLegal()
	assert.True(t, result, "Expected all injectors to be legal")

	// Test case 2: Injector with non-existent provider
	parser = &Parser{
		Packages: []*DiPackage{{
			Path: "github.com/my/package",
			Funcs: []*DiFunc{
				{
					Injectors: []*Injector{{
						ProviderId: "nonExistentProvider",
						Param:      "param1",
					}},
				},
			},
		}},
	}

	result = parser.checkInjectorLegal()
	assert.False(t, result, "Expected injectors to not be legal")
}

func TestParser_IncreaseProviderPrioritys(t *testing.T) {

	func1 := &DiFunc{
		ProviderId: "provider1",
	}
	func2 := &DiFunc{
		ProviderId: "provider2",
		Injectors: []*Injector{{
			ProviderId: "provider1",
			Dependency: func1,
		}},
	}

	parser := &Parser{
		Packages: []*DiPackage{{
			Funcs: []*DiFunc{func1, func2},
		}},
	}

	// Test case 1: No circular dependency
	result := parser.increaseProviderPrioritys(newChain(), parser.Packages[0].Funcs[0])
	assert.True(t, result, "Expected no circular dependency")

	// Test case 2: Circular dependency
	func1.Injectors = append(func1.Injectors, &Injector{
		ProviderId: "provider2",
		Dependency: func2,
	})
	func2.Injectors[0].Dependency = func1
	result = parser.increaseProviderPrioritys(newChain(), parser.Packages[0].Funcs[0])
	assert.False(t, result, "Expected circular dependency")
}

func TestParser_CheckCyclicProvider(t *testing.T) {

	func1 := &DiFunc{
		ProviderId: "provider1",
	}
	func2 := &DiFunc{
		ProviderId: "provider2",
		Injectors: []*Injector{{
			ProviderId: "provider1",
			Dependency: func1,
		}},
	}

	// Test case 1: No circular dependency
	parser := &Parser{
		Packages: []*DiPackage{{
			Funcs: []*DiFunc{func1, func2},
		}},
	}

	result := parser.checkCyclicProvider()
	assert.True(t, result, "Expected no circular dependency")
	assert.Equal(t, 1, func1.Sort)

	// Test case 2: Circular dependency
	func1.Injectors = append(func1.Injectors, &Injector{
		ProviderId: "provider2",
		Dependency: func2,
	})
	func2.Injectors[0].Dependency = func1

	parser = &Parser{
		Packages: []*DiPackage{{
			Funcs: []*DiFunc{func1, func2},
		}},
	}
	result = parser.checkCyclicProvider()
	assert.False(t, result, "Expected circular dependency")
}
