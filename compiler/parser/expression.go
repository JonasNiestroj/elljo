package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/parser"
)

func ReadExpression(src string) *ast.Program {
	parserInstance := parser.NewParser(src, 0)
	program := parserInstance.Parse()
	if len(parserInstance.Errors) > 0 {
		return nil
	}
	return program
}
