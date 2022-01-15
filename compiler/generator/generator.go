package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/sourcemap"
	"elljo/compiler/utils"
	"strings"
)

type FragmentCounter struct {
	Element int
	Text    int
	Anchor  int
}

type Fragment struct {
	Name               string
	UseAnchor          bool
	InitStatements     []Statement
	UpdateStatments    []Statement
	TeardownStatements []Statement
	ContextChain       []string
	Target             string
	Counters           FragmentCounter
	Parent             *Fragment
	IsComponent        bool
	HasContext         bool
	Mappings           [][]int
	UpdateContextChain string
	NeedsAppend        bool
}

type Statement struct {
	source   string
	mappings [][]int
}

type Generator struct {
	ifCounter           int
	elseCounter         int
	elseIfCounter       int
	eachCounter         int
	textCounter         int
	componentCounter    int
	renderers           []Renderer
	componentProperties []ComponentProperties
	FileName            string
	elements            map[string]struct{}
	loops               []string
}

type GeneratorOutput struct {
	Output    string `json:"output"`
	Sourcemap string `json:"sourcemap"`
	Css       string `json:"css"`
}

func (self *Generator) Visit(parser parser.Parser, children parser.Entry, current *Fragment, template string) {
	switch children.EntryType {
	case "Element":
		current = self.VisitElement(parser, children, current)
	case "Text":
		current = self.VisitText(children, current)
	case "MustacheTag":
		current = self.VisitMustache(parser, children, current)
	case "IfBlock":
		current = self.VisitIf(children, current)
	case "Loop":
		current = self.VisitLoop(children, current)
	case "ElseBlock":
		current = self.VisitElse(children, current)
	case "ElseIfBlock":
		current = self.VisitElseIf(children, current)
	}

	if len(children.Children) > 0 {
		for _, child := range children.Children {
			self.Visit(parser, *child, current, template)
		}
	}

	switch children.EntryType {
	case "Element":
		self.VisitElementAfter(current)
	case "IfBlock":
		self.VisitIfAfter(current)
	case "Loop":
		self.VisitLoopAfter(current)
	case "ElseBlock":
		self.VisitElseAfter(current)
	case "ElseIfBlock":
		self.VisitElseIfAfter(current)
	}
}

func (self *Generator) getJsSourceMap(parser parser.Parser, js string) [][]int {
	linesTillJs := utils.CountLines(parser.Template[0:parser.ScriptSource.StartIndex])

	jsLines := strings.Split(js, "\n")

	mappings := [][]int{{}, {}}

	for index, _ := range jsLines {
		if index == len(jsLines)-1 {
			break
		}
		lineIndex := linesTillJs + (index + 1)
		mappings = append(mappings, []int{0, 0, lineIndex, 0})
	}

	return mappings
}

func (self *Generator) getSourceMapMapping(parser parser.Parser, mappings [][]int) string {

	var sourceMapMapping strings.Builder
	var lastLine []int

	// Add ; for every import
	sourceMapMapping.WriteString(strings.Repeat(";", len(parser.ScriptSource.Imports)))

	for i := 0; i < len(mappings); i++ {
		mapping := mappings[i]
		if len(mapping) > 0 {
			if len(lastLine) != 0 {
				cpy := mapping[2]
				mapping[2] = mapping[2] - lastLine[2]
				lastLine = []int{mapping[0], mapping[1], cpy, mapping[3]}
			} else {
				lastLine = mapping
			}
			sourceMapMapping.WriteString(sourcemap.EncodeValues([]int{mapping[0], mapping[1], mapping[2], mapping[3]}) + ";")
		} else {
			sourceMapMapping.WriteString(";")
		}
	}

	return sourceMapMapping.String()
}

func (self *Generator) Generate(parser parser.Parser, template string) GeneratorOutput {
	self.elements = make(map[string]struct{})

	current := Fragment{
		UseAnchor:          false,
		Name:               "render",
		Target:             "target",
		InitStatements:     []Statement{},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		ContextChain:       []string{"context", "dirtyInState", "oldState"},
	}

	js := parser.ScriptSource.StringReplacer.String()

	decodedSourceMapMappings := self.getJsSourceMap(parser, js)

	self.Visit(parser, *parser.Entries[0], &current, template)

	self.renderers = append(self.renderers, self.CreateRenderer(current))

	utils.ReverseSlice(self.renderers)

	var renderersSources []string

	for _, renderer := range self.renderers {
		renderersSources = append(renderersSources, renderer.source)
		decodedSourceMapMappings = append(decodedSourceMapMappings, renderer.mappings...)
	}

	var code strings.Builder

	for _, importVar := range parser.ScriptSource.Imports {
		code.WriteString(importVar.Source + "\n")
	}

	sourceMapMapping := self.getSourceMapMapping(parser, decodedSourceMapMappings)

	properties := ""

	if len(parser.ScriptSource.Properties) > 0 {
		properties = "let that = this;"
	}

	for _, property := range parser.ScriptSource.Properties {
		properties += `
			Object.defineProperty(this.$props, "` + property + `", {
				get() {
					return this.` + property + `;
				},
				set(value) {
					that.updateValue("` + property + `", ` + property + ` = value);
				}
			})
			if(props['` + property + `']) {
				` + property + ` = props['` + property + `'];
			}`
	}

	elementCache := "const elementCache = {};"

	for key, _ := range self.elements {
		elementCache += `
			elementCache.` + key + ` = document.createElement("` + key + `");`
	}

	code.WriteString(`import { setComponent, EllJoComponent, createFragment } from '@elljo/runtime'
		` + elementCache + `
		class ` + self.FileName + ` extends EllJoComponent {
			constructor(options, props, events) {
				super(options, props, events)
				this.init(options, props, events);
			}

			update() {
				super.update()
			}

			init(options, props, events) {
				` + js + strings.Join(renderersSources, "\n") + `

				` + properties + `		

				this.$.mainFragment = render(options.target);
				this.queueUpdate();
			}
		}
		export default ` + self.FileName)

	stringReplacer := utils.StringReplacer{Text: parser.StyleSource.Source}

	for _, rule := range parser.StyleSource.Rules {
		stringReplacer.Replace(rule.StartIndex, rule.EndIndex, rule.Selector+"[scope-"+parser.StyleSource.Id+"]")
	}

	return GeneratorOutput{Output: code.String(), Sourcemap: sourcemap.CreateSourcemap(parser.FileName, sourceMapMapping), Css: stringReplacer.String()}
}
