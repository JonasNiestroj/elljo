package ssr

import (
	"elljo/compiler/parser"
)

import v8 "rogchap.com/v8go"

type SSR struct {
	Render string
}

type Output struct {
	Html string `json:"html"`
	Css  string `json:"css"`
}

func (self *SSR) Visit(parser parser.Parser, children parser.Entry) {
	switch children.EntryType {
	case "Element":
		self.VisitElement(parser, children)
	case "Text":
		self.VisitText(parser, children)
	case "MustacheTag":
		self.VisitMustache(parser, children)
	}

	if len(children.Children) > 0 {
		for _, child := range children.Children {
			self.Visit(parser, *child)
		}
	}

	switch children.EntryType {
	case "Element":
		self.VisitElementAfter(children)
	}
}

func (self *SSR) SSR(parser parser.Parser) Output {

	self.Visit(parser, *parser.Entries[0])

	style := parser.StyleSource.Source
	indexToAdd := 0
	for _, rule := range parser.StyleSource.Rules {
		style = style[0:rule.StartIndex+indexToAdd] + rule.Selector + "[scope-" + parser.StyleSource.Id + "]" + style[rule.EndIndex+indexToAdd:]
		indexToAdd += 8 + len(parser.StyleSource.Id)
	}
	code := `
var Main = {};`

	code += `
Main.data = {
`

	for _, variable := range parser.ScriptSource.Variables {
		code += `
        ` + variable.Name + `: ` + parser.ScriptSource.Source[variable.Initializer.Index0():variable.Initializer.Index1()]
	}

	code += `
}`

	code += `
Main.render = function() {
 return ` + "`" + self.Render + "`" + `
}`

	ctx, _ := v8.NewContext()
	ctx.RunScript(code, "code.js")
	ctx.RunScript("const rendered = Main.render()", "render.js")
	renderedHtml, _ := ctx.RunScript("rendered", "rendered.js")

	return Output{Html: renderedHtml.String(), Css: style}
}
