package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/parser"
	"elljo/compiler/utils"
	"regexp"
	"strings"
)

type Import struct {
	Source  string
	Name    string
	IsEllJo bool
}

type Property struct {
	ExportStatement *ast.ExportStatement
	Name            string
}

type Variable struct {
	Name        string
	Initializer ast.Expression
	IsProperty  bool
}

type ParsedVariable struct {
	Name       string
	Value      string
	IsProperty string
}

type Assign struct {
	Name       string
	Expression *ast.AssignExpression
}

var (
	variables     []Variable
	thisVariables []string
	imports       []*ast.ImportStatement
	dotExpression *ast.DotExpression
	properties    []Property
	assigns       []Assign
)

func Walk(node ast.Node) {
	if id, ok := node.(*ast.ImportStatement); ok && id != nil {
		imports = append(imports, id)
	}
	if id, ok := node.(*ast.ThisExpression); ok && id != nil {
		if dotExpression != nil {
			thisVariables = append(thisVariables, dotExpression.Identifier.Name.String())
		}
	}

	dotExpression = nil
	if id, ok := node.(*ast.DotExpression); ok && id != nil {
		for _, variable := range thisVariables {
			if variable == id.Identifier.Name.String() {
				return
			}
		}
		dotExpression = id
	}

	if id, ok := node.(*ast.VariableStatement); ok && id != nil {
		for _, variable := range id.List {
			if id, ok := variable.(*ast.VariableExpression); ok && id != nil {
				variables = append(variables, Variable{
					Name:        id.Name.String(),
					Initializer: id.Initializer,
				})
			}
		}
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

	if expression, ok := node.(*ast.AssignExpression); ok && expression != nil {
		if left, ok := expression.Left.(*ast.Identifier); ok && left != nil {
			if variablesContains(left.Name.String()) {
				assigns = append(assigns, Assign{
					Name:       left.Name.String(),
					Expression: expression,
				})
			}
		}
		if left, ok := expression.Left.(*ast.DotExpression); ok && left != nil {
			if variablesContains(left.Identifier.Name.String()) {
				assigns = append(assigns, Assign{
					Name:       left.Identifier.Name.String(),
					Expression: expression,
				})
			}
		}
	}
}

func variablesContains(name string) bool {
	for _, variable := range variables {
		if variable.Name == name {
			return true
		}
	}
	return false
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

	stringReplacer := &utils.StringReplacer{
		Text: source,
	}

	var importNames []Import

	for _, importStatement := range imports {
		isElljoFile := false
		// TODO: Improve .jo check
		if !strings.HasSuffix(importStatement.Source, ".jo'") && !strings.HasSuffix(importStatement.Source, ".jo\"") {
			isElljoFile = true
		}
		importVar := Import{
			Source:  source[importStatement.Index0():importStatement.Index1()],
			Name:    importStatement.Name,
			IsEllJo: isElljoFile,
		}
		stringReplacer.Replace(importStatement.Index0(), importStatement.Index1(), "")
		importNames = append(importNames, importVar)
	}

	for _, assign := range assigns {
		assignSource := source[assign.Expression.Index0():assign.Expression.Index1()]
		newAssignSource := "this.oldState." + assign.Name + " = " + assign.Name + "; this.updateValue('" + assign.Name + "', " + assignSource + ");"
		stringReplacer.Replace(assign.Expression.Index0(), assign.Expression.Index1(), newAssignSource)
	}

	var propertyNames []string
	for _, export := range properties {
		exportStatement := export.ExportStatement
		stringReplacer.Replace(exportStatement.Export, exportStatement.Statement.Index0(), "")
		propertyNames = append(propertyNames, export.Name)
		variables = append(variables, Variable{
			Name:       export.Name,
			IsProperty: true,
		})
	}

	// var usedVariables []ParsedVariable

	/*for _, thisVariable := range thisVariables {
		contains := false
		for _, variable := range variables {
			if variable.Name == thisVariable {
				contains = true
			}
		}

		if contains {
			usedVariables = append(usedVariables, thisVariable)
		}
	}*/

	/*for _, properties := range properties {
		usedVariables = append(usedVariables, properties.Name)
	}*/

	end := parserInstance.Index
	parserInstance.Index += 9

	return ScriptSource{
		StartIndex:     start,
		EndIndex:       end,
		Program:        program,
		Variables:      variables,
		Imports:        importNames,
		Source:         source,
		Properties:     propertyNames,
		StringReplacer: stringReplacer,
	}
}
