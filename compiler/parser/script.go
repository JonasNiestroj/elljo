package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/parser"
	"regexp"
	"strings"
)

var (
	variables []string
	assigns []*ast.AssignExpression
	imports []*ast.ImportStatement
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
	if id, ok := node.(*ast.AssignExpression); ok && id != nil {
		assigns = append(assigns, id)
	}
	if id, ok := node.(*ast.ImportStatement); ok && id != nil {
		imports = append(imports, id)
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
		if id, ok := declaration.(*ast.VariableDeclaration); ok && id != nil {
			for _, item := range id.List {
				variables = append(variables, item.Name.String())
			}
		}
		if id, ok := declaration.(*ast.FunctionDeclaration); ok && id != nil {
			ast.Walk(id.Function, Walk)
		}
	}

	indexToAdd := 0

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

	for _, assignExpression := range assigns {
		println(source[assignExpression.Index0():assignExpression.Index1()])
		if identifier, ok := assignExpression.Left.(*ast.Identifier); ok && identifier != nil {
			parserInstance.Template = parserInstance.Template[:assignExpression.Index1() + indexToAdd] + `;currentComponent.set({` + identifier.Name.String() + `}, "`+ identifier.Name.String() + `");` + parserInstance.Template[int(assignExpression.Index1()) + indexToAdd:]
			indexToAdd += 30 + len(identifier.Name) * 2
			parserInstance.Index += 30 + len(identifier.Name) * 2
		}
	}

	return ScriptSource{
		StartIndex: start,
		EndIndex:   parserInstance.Index,
		Program:    program,
		Variables:  variables,
		Imports: importNames,
	}
}
