package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/parser"
)

func ReadExpression(src string) *ast.Program {
	parserInstance := parser.NewParser(src, 0)
	program, err := parserInstance.Parse()
	if err != nil {
		panic(err)
	}
	return program
}
