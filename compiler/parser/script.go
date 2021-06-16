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

type Property struct {
	ExportStatement *ast.ExportStatement
	Name            string
}

var (
	variables     []string
	imports       []*ast.ImportStatement
	dotExpression *ast.DotExpression
	properties    []Property
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

	if expression, ok := node.(*ast.ExportStatement); ok && expression != nil {
		if id, ok := expression.Statement.(*ast.VariableStatement); ok && id != nil {
			for _, variable := range id.List {
				if id, ok := variable.(*ast.VariableExpression); ok && id != nil {
					properties = append(properties, Property{
						ExportStatement: expression,
						Name:            id.Name.String(),
					})
				}
			}
		}
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

	var propertyNames []string
	for _, export := range properties {
		exportStatement := export.ExportStatement
		source = source[:exportStatement.Export-indexToRemove] + source[exportStatement.Statement.Index0()-indexToRemove:]
		indexToRemove += 6
		propertyNames = append(propertyNames, export.Name)
		variables = append(variables, export.Name)
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
		Properties: propertyNames,
	}
}
