package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/sourcemap"
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
}

type Statement struct {
	source   string
	mappings [][]int
}

type Generator struct {
	ifCounter     int
	elseCounter   int
	elseIfCounter int
	eachCounter   int
	textCounter   int
	renderers     []Renderer
}

type GeneratorOutput struct {
	Output    string `json:"output"`
	Sourcemap string `json:"sourcemap"`
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
	current := Fragment{
		UseAnchor:          false,
		Name:               "render",
		Target:             "target",
		InitStatements:     []Statement{},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		ContextChain:       []string{"context", "dirtyInState", "oldState"},
	}

	js := template[parser.ScriptSource.StartIndex:parser.ScriptSource.EndIndex]

	linesTillJs := 0

	for _, c := range template[0:parser.ScriptSource.StartIndex] {
		if c == '\n' {
			linesTillJs++
		}
	}

	mappings := [][]int{{}}

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

	code := `var component = function(options) {` + js +
		`; var currentComponent = null;
		` + strings.Join(renderersSources, "\n") +
		`
			var component = {};
			var state = {};
			var oldState = {};
			var updating = false;
			var dirtyInState = [];
			var update = function() {
				if(!updating) {
					updating = true;
					Promise.resolve().then(component.update)
				}
			}
			component.set = function set (newState, name) {
				dirtyInState.push(name);
				if(!oldState) {
					oldState = Object.assign({}, state)
				} else {
					oldState[name] = state[name];
				}
				Object.assign(state, newState);
				update()
			};
			component.update = function update() {
				updating = false;
				mainFragment.update(state, dirtyInState, oldState);
				dirtyInState = [];
				oldState = {}
			}
			component.teardown = function teardown () {
				mainFragment.teardown();
				mainFragment = null;
				state = {};
			};
			this.contexts = [];
        	this.utils = { diffArray: function diffArray(one, two) {
                if (!Array.isArray(two)) {
                    return one.slice();
                }

                var tlen = two.length
                var olen = one.length;
                var idx = -1;
                var arr = [];

                while (++idx < olen) {
                    var ele = one[idx];

                    var hasEle = false;
                    for (var i = 0; i < tlen; i++) {
                        var val = two[i];

                        if (ele === val) {
                            hasEle = true;
                            break;
                        }
                    }

                    if (hasEle === false) {
                        arr.push({element: ele, index: idx});
                    }
                }
                return arr;
            } }
			let mainFragment = render(options.target);
			component.set({` + strings.Join(parser.ScriptSource.Variables, ",") + `});
			currentComponent = component 
			return component;
	}
	export default component`

	return GeneratorOutput{Output: code, Sourcemap: sourcemap.CreateSourcemap(mappingsStrings)}
}
