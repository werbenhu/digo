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
)

func TestNewIdent(t *testing.T) {
	ident := newIdent("name")
	assert.Equal(t, "name", ident.Name)
}

func TestNewSelectorExpr(t *testing.T) {
	selectorExpr := newSelectorExpr("pkg.name")
	splitted := strings.Split("pkg.name", ".")
	assert.Equal(t, splitted[0], selectorExpr.X.(*ast.Ident).Name)
	assert.Equal(t, splitted[1], selectorExpr.Sel.Name)
}

func TestNewStarExpr(t *testing.T) {
	starExpr := newStarExpr("pkg.name")
	splitted := strings.Split("pkg.name", ".")
	selectorExpr := starExpr.X.(*ast.SelectorExpr)
	assert.Equal(t, splitted[0], selectorExpr.X.(*ast.Ident).Name)
	assert.Equal(t, splitted[1], selectorExpr.Sel.Name)
}

func TestNewCommentGroup(t *testing.T) {
	texts := []string{"comment 1", "comment 2"}
	commentGroup := newCommentGroup(texts)
	assert.Len(t, commentGroup.List, len(texts))
	for i, comment := range commentGroup.List {
		assert.Equal(t, texts[i], comment.Text)
	}
}

func TestNewCallExpr(t *testing.T) {
	fn := newIdent("fn")
	args := []ast.Expr{newIdent("arg1"), newIdent("arg2")}
	callExpr := newCallExpr(fn, args)
	assert.Equal(t, fn, callExpr.Fun)
	assert.Equal(t, args, callExpr.Args)
}

func TestNewExprs(t *testing.T) {
	exprs := []ast.Expr{newIdent("expr1"), newIdent("expr2")}
	newExprs := newExprs(exprs...)
	assert.Len(t, newExprs, len(exprs))
	for i, expr := range newExprs {
		assert.Equal(t, exprs[i], expr)
	}
}

func TestNewBasicLit(t *testing.T) {
	val := "value"
	basicLit := newBasicLit(val)
	assert.Equal(t, token.STRING, basicLit.Kind)
	assert.Equal(t, "\""+val+"\"", basicLit.Value)
}

func TestNewErrCheckStmt(t *testing.T) {
	errCheckStmt := newErrCheckStmt()
	ifStmt, ok := errCheckStmt.(*ast.IfStmt)
	assert.True(t, ok)
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	assert.True(t, ok)
	assert.Equal(t, newIdent("err"), binaryExpr.X)
	assert.Equal(t, token.NEQ, binaryExpr.Op)
	assert.Equal(t, newIdent("nil"), binaryExpr.Y)
	assert.NotNil(t, ifStmt.Body)
	assert.Len(t, ifStmt.Body.List, 1)
	exprStmt, ok := ifStmt.Body.List[0].(*ast.ExprStmt)
	assert.True(t, ok)
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	assert.True(t, ok)
	assert.Equal(t, newIdent("panic"), callExpr.Fun)
	assert.NotNil(t, callExpr.Args)
	assert.Len(t, callExpr.Args, 1)
	assert.Equal(t, newIdent("err"), callExpr.Args[0])
}

func TestNewImportSpec(t *testing.T) {
	path := "github.com/example/pkg"
	alias := "pkgalias"
	importSpec := newImportSpec(path, alias)
	assert.Equal(t, newBasicLit(path), importSpec.Path)
	assert.NotNil(t, importSpec.Name)
	assert.Equal(t, alias, importSpec.Name.Name)
}

func TestObjName(t *testing.T) {
	testCases := []struct {
		prefix string
		expect string
	}{
		{
			prefix: "pkg.name",
			expect: "pkg_name_obj",
		},
		{
			prefix: "path/to/package",
			expect: "path_to_package_obj",
		},
		{
			prefix: "prefix",
			expect: "prefix_obj",
		},
		{
			prefix: "prefix.with.multiple.dots",
			expect: "prefix_with_multiple_dots_obj",
		},
	}

	for _, tc := range testCases {
		name := objName(tc.prefix)
		assert.Equal(t, tc.expect, name)
	}
}

func TestDefineInjectStmts(t *testing.T) {

	g := NewGenerator(nil)

	inject := &Injector{
		Param:      "myParam",
		Pkg:        "github.com/werbenhu/eventbus",
		Alias:      "bus",
		ProviderId: "my_provider_id",
		Typ:        newIdent("MyType"),
	}

	// Define the inject statements.
	stmts := g.defineInjectStmts(inject)

	// Assert the number of statements.
	assert.Len(t, stmts, 3)

	// Assert ImportSpecs
	assert.Len(t, g.ImportSpecs, 1)
	importSpec, ok := g.ImportSpecs["github.com/werbenhu/eventbus_bus"]
	assert.True(t, ok)
	assert.Equal(t, &ast.ImportSpec{
		Path: newBasicLit("github.com/werbenhu/eventbus"),
		Name: newIdent("bus"),
	}, importSpec)

	// Assert the assignment statement for providing the object.
	assert.Equal(t, &ast.AssignStmt{
		Lhs: newExprs(newIdent(objName("myParam")), newIdent("err")),
		Tok: token.DEFINE,
		Rhs: newExprs(
			newCallExpr(
				newSelectorExpr(g.ProvideFunction),
				[]ast.Expr{newBasicLit("my_provider_id")},
			),
		),
	}, stmts[0])

	// Assert the error check statement.
	assert.Equal(t, newErrCheckStmt(), stmts[1])

	// Assert the assignment statement for type assertion.
	assert.Equal(t, &ast.AssignStmt{
		Lhs: newExprs(newIdent("myParam")),
		Tok: token.DEFINE,
		Rhs: newExprs(&ast.TypeAssertExpr{
			X:    newIdent(objName("myParam")),
			Type: newIdent("MyType"),
		}),
	}, stmts[2])

}
