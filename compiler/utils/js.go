package utils

import (
	"github.com/JonasNiestroj/esbuild-internal/js_ast"
)

func IsOnlyStringExpression(jsAst js_ast.AST) bool {
	isStringExpression := false

	for _, part := range jsAst.Parts {
		for _, stmt := range part.Stmts {
			if _, ok := stmt.Data.(*js_ast.SDirective); ok {
				isStringExpression = true
			}
		}
	}

	return isStringExpression
}
