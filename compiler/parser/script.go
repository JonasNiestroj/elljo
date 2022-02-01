package parser

import (
	"elljo/compiler/utils"
	"github.com/JonasNiestroj/esbuild-internal/config"
	"github.com/JonasNiestroj/esbuild-internal/js_ast"
	"github.com/JonasNiestroj/esbuild-internal/js_parser"
	"github.com/JonasNiestroj/esbuild-internal/logger"
	"regexp"
	"strings"
)

type Import struct {
	Source  string
	Name    string
	IsEllJo bool
}

type Property struct {
	ExportStatement js_ast.NamedExport
	Name            string
}

type Variable struct {
	Name       string
	IsProperty bool
}

type ParsedVariable struct {
	Name       string
	Value      string
	IsProperty string
}

type Assign struct {
	Name       string
	Expression *js_ast.EBinary
}

var (
	variables  []Variable
	imports    []*js_ast.SImport
	properties []Property
	assigns    []Assign
	symbols    js_ast.SymbolMap
	astTree    js_ast.AST
)

func WalkFunc(entry interface{}) {
	if expr, ok := entry.(*js_ast.EBinary); ok && expr != nil {
		if left, ok := expr.Left.Data.(*js_ast.EIdentifier); ok && left != nil {
			variableName := symbols.Get(left.Ref).OriginalName
			if variablesContains(variableName) {
				assigns = append(assigns, Assign{
					Name:       variableName,
					Expression: expr,
				})
			}
		}
	}

	if imp, ok := entry.(*js_ast.SImport); ok && imp != nil {
		imports = append(imports, imp)
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

	log := logger.NewDeferLog(logger.DeferLogAll)

	src := logger.Source{Index: 0, Contents: source}

	success := false

	astTree, success = js_parser.Parse(log, src, js_parser.OptionsFromConfig(&config.Options{WriteToStdout: false}))

	if !success {
		msgs := log.Done()

		var errors []utils.Error

		for _, msg := range msgs {
			err := utils.Error{
				Line:        msg.Data.Location.Line,
				Message:     msg.Data.Text,
				StartColumn: msg.Data.Location.Column,
				EndColumn:   msg.Data.Location.Column + msg.Data.Location.Length,
			}

			errors = append(errors, err)
		}

		parserInstance.Errors = append(parserInstance.Errors, errors...)

		return ScriptSource{}
	}

	symbols = js_ast.NewSymbolMap(1)
	symbols.SymbolsForSource[0] = astTree.Symbols

	for _, part := range astTree.Parts {
		for _, symbol := range part.DeclaredSymbols {
			if symbol.IsTopLevel {
				if symbols.Get(symbol.Ref).Kind != js_ast.SymbolImport {
					variables = append(variables, Variable{
						Name: symbols.Get(symbol.Ref).OriginalName,
					})
				}
			}
		}
		for _, stmt := range part.Stmts {
			Walk(stmt.Data, WalkFunc)
		}
	}

	// Get exports for properties
	for key, value := range astTree.NamedExports {
		properties = append(properties, Property{
			ExportStatement: value,
			Name:            key,
		})
	}

	stringReplacer := &utils.StringReplacer{
		Text: source,
	}

	var importNames []Import

	for _, importStatement := range imports {
		index0 := astTree.NamedImports[importStatement.NamespaceRef].AliasLoc.Start
		index1 := astTree.ImportRecords[importStatement.ImportRecordIndex].Range.End()

		importFile := astTree.ImportRecords[importStatement.ImportRecordIndex].Path.Text

		name := ""

		if importStatement.DefaultName != nil {
			name = symbols.Get(importStatement.DefaultName.Ref).OriginalName
		}

		isElljoFile := false
		// TODO: Improve .jo check
		if !strings.HasSuffix(importFile, ".jo'") && !strings.HasSuffix(importFile, ".jo\"") {
			isElljoFile = true
		}

		importVar := Import{
			Source:  source[index0:index1],
			Name:    name,
			IsEllJo: isElljoFile,
		}

		importNames = append(importNames, importVar)

		stringReplacer.Replace(int(index0), int(index1), "")
	}

	for _, assign := range assigns {
		if left, ok := assign.Expression.Left.Data.(*js_ast.EIdentifier); ok && left != nil {
			endNewLine := src.RangeOfOperatorAfter(assign.Expression.Right.Loc, "\n")
			endSemicolon := src.RangeOfOperatorAfter(assign.Expression.Right.Loc, ";")

			var end logger.Range

			if endNewLine.End() < endSemicolon.End() {
				end = endNewLine
			} else {
				end = endSemicolon
			}

			index0 := int(assign.Expression.Left.Loc.Start)
			index1 := int(end.End() - 1)

			assignSource := source[index0:index1]
			newAssignSource := "this.oldState." + assign.Name + " = " + assign.Name + "; this.updateValue('" + assign.Name + "', " + assignSource + ");"

			stringReplacer.Replace(index0, index1, newAssignSource)
		}
	}

	var propertyNames []string
	for _, export := range properties {
		index0 := int(astTree.ExportKeyword.Loc.Start)
		index1 := int(astTree.ExportKeyword.End())

		stringReplacer.Replace(index0, index1, "")
		propertyNames = append(propertyNames, export.Name)
	}

	end := parserInstance.Index
	parserInstance.Index += 9

	return ScriptSource{
		StartIndex:     start,
		EndIndex:       end,
		Variables:      variables,
		Imports:        importNames,
		Source:         source,
		Properties:     propertyNames,
		StringReplacer: stringReplacer,
	}
}
