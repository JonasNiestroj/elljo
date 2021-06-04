package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/parser"
	"regexp"
	"strings"
)

var (
	variables     []string
	imports       []*ast.ImportStatement
	dotExpression *ast.DotExpression
)

func Spaces(count int) string {
	result := ""
	for count != 0 {
		result += " "
		count--
	}
	return result
}

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

	source := Spaces(scriptStart) + parserInstance.Template[scriptStart:parserInstance.Index]

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

	var importNames []string

	for _, importStatement := range imports {
		// TODO: Improve .jo check
		if !strings.HasSuffix(importStatement.Source, ".jo'") && !strings.HasSuffix(importStatement.Source, ".jo\"") {
			continue
		}
		/*parserInstance.Template = parserInstance.Template[:importStatement.Index0()] + parserInstance.Template[importStatement.Index1():]
		indexToAdd -= int(importStatement.Index1()) - int(importStatement.Index0())
		parserInstance.Index -= int(importStatement.Index1()) - int(importStatement.Index0())*/
		importNames = append(importNames, importStatement.Name)
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
