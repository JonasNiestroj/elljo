package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/parser"
	"regexp"
	"strings"
)

type Import struct {
	Source string
	Name   string
}

var (
	variables     []string
	imports       []*ast.ImportStatement
	dotExpression *ast.DotExpression
)

func Walk(node ast.Node) {
	if id, ok := node.(*ast.ImportStatement); ok && id != nil {
		imports = append(imports, id)
	}
	if id, ok := node.(*ast.ThisExpression); ok && id != nil {
		if dotExpression != nil {
			variables = append(variables, dotExpression.Identifier.Name.String())
		}
	}
	dotExpression = nil
	if id, ok := node.(*ast.DotExpression); ok && id != nil {
		for _, variable := range variables {
			if variable == id.Identifier.Name.String() {
				return
			}
		}
		dotExpression = id
	}
}

func ReadScript(parserInstance *Parser, start int) ScriptSource {
	scriptStart := parserInstance.Index
	pattern, _ := regexp.Compile("</script>")
	parserInstance.ReadUntil(pattern)

	source := parserInstance.Template[scriptStart:parserInstance.Index]

	jsParserInstance := parser.NewParser(source, 0)

	program, err := jsParserInstance.Parse()

	if err != nil {
		panic(err)
	}

	ast.Walk(program, Walk)

	for _, declaration := range program.DeclarationList {

		if id, ok := declaration.(*ast.FunctionDeclaration); ok && id != nil {
			ast.Walk(id.Function, Walk)
		}
	}

	var importNames []Import
	indexToRemove := 0
	for _, importStatement := range imports {
		// TODO: Improve .jo check
		if !strings.HasSuffix(importStatement.Source, ".jo'") && !strings.HasSuffix(importStatement.Source, ".jo\"") {
			continue
		}
		importVar := Import{
			Source: source[importStatement.Index0()-indexToRemove : importStatement.Index1()-indexToRemove],
			Name:   importStatement.Name,
		}
		source = source[:importStatement.Index0()-indexToRemove] +
			source[importStatement.Index1()-indexToRemove:]
		indexToRemove += int(importStatement.Index1()) - int(importStatement.Index0())
		importNames = append(importNames, importVar)
	}

	end := parserInstance.Index
	parserInstance.Index += 9
	return ScriptSource{
		StartIndex: start,
		EndIndex:   end,
		Program:    program,
		Variables:  variables,
		Imports:    importNames,
		Source:     source,
	}
}
