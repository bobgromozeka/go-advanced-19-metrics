package osexitcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "OsExitCheck",
	Run:  run,
	Doc:  "Check for os.Exit calls in main function",
}

func run(pass *analysis.Pass) (any, error) {
	if pass.Pkg.Name() == "main" && !strings.Contains(pass.Pkg.Path(), ".test") { // Checking only main packages
		for _, f := range pass.Files {
			ast.Inspect(
				f, func(node ast.Node) bool {
					if fd, ok := node.(*ast.FuncDecl); ok && fd.Name.Name == "main" { // main function
						for _, stmt := range fd.Body.List {
							if expr, ok := stmt.(*ast.ExprStmt); ok {
								if call, ok := expr.X.(*ast.CallExpr); ok {
									if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
										xIdent, xOk := fun.X.(*ast.Ident)
										if xOk && xIdent.Name == "os" && fun.Sel.Name == "Exit" {
											pass.Reportf(fun.Pos(), "os.Exit call in main function")
										}
									}
								}
							}
						}
					}
					return true
				},
			)
		}
	}
	return nil, nil
}
