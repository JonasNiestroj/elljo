package parser

import (
	"github.com/JonasNiestroj/esbuild-internal/js_ast"
)

func Walk(entry interface{}, nodeCallback func(entry interface{})) {

	switch entry := entry.(type) {
	case *js_ast.SExpr:
		nodeCallback(entry)
		Walk(entry.Value.Data, nodeCallback)
	case *js_ast.SFunction:
		nodeCallback(entry)
		for _, stmt := range entry.Fn.Body.Stmts {
			Walk(stmt.Data, nodeCallback)
		}
	case *js_ast.SBlock:
		nodeCallback(entry)
		for _, stmt := range entry.Stmts {
			Walk(stmt.Data, nodeCallback)
		}
	case *js_ast.SDoWhile:
		nodeCallback(entry)
		Walk(entry.Body, nodeCallback)
	case *js_ast.SFor:
		nodeCallback(entry)
		Walk(entry.InitOrNil, nodeCallback)
		Walk(entry.TestOrNil, nodeCallback)
		Walk(entry.UpdateOrNil, nodeCallback)
		Walk(entry.Body, nodeCallback)
	case *js_ast.SForIn:
		nodeCallback(entry)
		Walk(entry.Body, nodeCallback)
	case *js_ast.SForOf:
		nodeCallback(entry)
		Walk(entry.Body, nodeCallback)
	case *js_ast.STry:
		nodeCallback(entry)
		for _, stmt := range entry.Body {
			Walk(stmt, nodeCallback)
		}
		for _, stmt := range entry.Catch.Body {
			Walk(stmt, nodeCallback)
		}
	case *js_ast.SWhile:
		nodeCallback(entry)
		Walk(entry.Test, nodeCallback)
		Walk(entry.Body, nodeCallback)
	case *js_ast.SWith:
		nodeCallback(entry)
		Walk(entry.Value, nodeCallback)
		Walk(entry.Body, nodeCallback)
	case *js_ast.SImport:
		nodeCallback(entry)
	case *js_ast.SLocal:
		nodeCallback(entry)
		for _, decl := range entry.Decls {
			Walk(decl.ValueOrNil.Data, nodeCallback)
		}
	case *js_ast.SIf:
		nodeCallback(entry)
		Walk(entry.Yes.Data, nodeCallback)
		Walk(entry.NoOrNil.Data, nodeCallback)
	case *js_ast.EBinary:
		nodeCallback(entry)
		Walk(entry.Left, nodeCallback)
		Walk(entry.Right, nodeCallback)
	case *js_ast.EArrow:
		nodeCallback(entry)
		for _, stmt := range entry.Body.Stmts {
			Walk(stmt.Data, nodeCallback)
		}
	case *js_ast.EFunction:
		nodeCallback(entry)
		for _, stmt := range entry.Fn.Body.Stmts {
			Walk(stmt, nodeCallback)
		}
	case *js_ast.EIf:
		nodeCallback(entry)
		Walk(entry.Yes.Data, nodeCallback)
		Walk(entry.No.Data, nodeCallback)
	case *js_ast.EIdentifier:
		nodeCallback(entry)
	case *js_ast.FnBody:
		nodeCallback(entry)
		for _, stmt := range entry.Stmts {
			Walk(stmt.Data, nodeCallback)
		}
	case *js_ast.Finally:
		nodeCallback(entry)
		for _, stmt := range entry.Stmts {
			Walk(stmt.Data, nodeCallback)
		}
	case *js_ast.ClassStaticBlock:
		nodeCallback(entry)
		for _, stmt := range entry.Stmts {
			Walk(stmt.Data, nodeCallback)
		}
	}
}
