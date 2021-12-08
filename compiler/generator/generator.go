package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/sourcemap"
	"strconv"
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

	linesTillJs := 0

	for _, c := range template[0:parser.ScriptSource.StartIndex] {
		if c == '\n' {
			linesTillJs++
		}
	}

	mappings := [][]int{{}, {}}

	for index, _ := range strings.Split(js, "\n") {
		if index == len(strings.Split(js, "\n"))-1 {
			break
		}
		lineIndex := linesTillJs + (index + 1)
		mappings = append(mappings, []int{0, 0, lineIndex, 0})
	}

	self.Visit(parser, *parser.Entries[0], &current, template)

	self.renderers = append(self.renderers, self.CreateRenderer(current))

	for i, j := 0, len(self.renderers)-1; i < j; i, j = i+1, j-1 {
		self.renderers[i], self.renderers[j] = self.renderers[j], self.renderers[i]
	}

	var renderersSources []string

	for _, renderer := range self.renderers {
		renderersSources = append(renderersSources, renderer.source)
		mappings = append(mappings, renderer.mappings...)
	}

	var mappingsStrings []string

	var lastLine []int

	code := ""

	for _, importVar := range parser.ScriptSource.Imports {
		code += importVar.Source + `
`
		mappingsStrings = append(mappingsStrings, ";")
	}

	for i := 0; i < len(mappings); i++ {
		mapping := mappings[i]
		if len(mapping) > 0 {
			if len(lastLine) != 0 {
				copy := mapping[2]
				mapping[2] = mapping[2] - lastLine[2]
				lastLine = []int{mapping[0], mapping[1], copy, mapping[3]}
			} else {
				lastLine = mapping
			}
			mappingsStrings = append(mappingsStrings, sourcemap.EncodeValues([]int{mapping[0], mapping[1], mapping[2], mapping[3]})+";")
		} else {
			mappingsStrings = append(mappingsStrings, ";")
		}
	}

	for _, variable := range parser.ScriptSource.Variables {
		propertyUpdate := ""
		for _, componentProperties := range self.componentProperties {
			if id, ok := componentProperties.Properties[variable.Name]; ok {
				propertyUpdate += `
					this['component-` + strconv.Itoa(componentProperties.Index) + `'].$props['` + id + `'] = value`
			}
		}

	}

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
					that.` + property + ` = value;
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

	code += `import { setComponent, EllJoComponent, Observer } from '@elljo/runtime'
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
		export default ` + self.FileName

	style := parser.StyleSource.Source
	indexToAdd := 0
	for _, rule := range parser.StyleSource.Rules {
		style = style[0:rule.StartIndex+indexToAdd] + rule.Selector + "[scope-" + parser.StyleSource.Id + "]" + style[rule.EndIndex+indexToAdd:]
		indexToAdd += 8 + len(parser.StyleSource.Id)
	}

	return GeneratorOutput{Output: code, Sourcemap: sourcemap.CreateSourcemap(parser.FileName, mappingsStrings), Css: style}
}
