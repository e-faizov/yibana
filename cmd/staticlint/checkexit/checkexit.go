package checkexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `check for os.Exit in main func main package`

var Analyzer = &analysis.Analyzer{
	Name:     "checkosexit",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	mainFunc := func(list []ast.Stmt) {
		for _, stmt := range list {
			switch x := stmt.(type) {
			case *ast.ExprStmt:
				// Поиск вызова os.Exit
				call, ok := x.X.(*ast.CallExpr)
				if ok {
					if identName, okName := call.Fun.(*ast.SelectorExpr); okName {
						if identPack, okIdent := identName.X.(*ast.Ident); okIdent {
							if identPack.Name == "os" && identName.Sel.Name == "Exit" {
								pass.Reportf(identName.Pos(), "Сall os.Exit in main function main package.")
							}
						}
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		// Поиск пакета main
		if file.Name.Name == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.FuncDecl:
					// Поиск функции main в пакете main
					if x.Name.Name == "main" {
						mainFunc(x.Body.List)
						return false
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
