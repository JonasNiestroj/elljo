package generator

import (
	"elljo/compiler/parser"
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
	InitStatements     []string
	UpdateStatments    []string
	TeardownStatements []string
	ContextChain       []string
	Target             string
	Counters           FragmentCounter
	Parent             *Fragment
	IsComponent        bool
	HasContext 		   bool
}

type Generator struct {
	ifCounter   int
	eachCounter int
	textCounter int
	renderers   []string
}

func (self *Generator) Visit(parser parser.Parser, children parser.Entry, current *Fragment, template string) {
	switch children.EntryType {
	case "Element":
		current = self.VisitElement(parser, children, current)
	case "Text":
		current = self.VisitText(children, current)
	case "MustacheTag":
		current = self.VisitMustache(children, current)
	case "IfBlock":
		current = self.VisitIf(children, current)
	case "Loop":
		current = self.VisitLoop(children, current)
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
	}
}

func (self *Generator) Generate(parser parser.Parser, template string) string {
	current := Fragment{
		UseAnchor:          false,
		Name:               "render",
		Target:             "target",
		InitStatements:     []string{},
		UpdateStatments:    []string{},
		TeardownStatements: []string{},
		ContextChain:       []string{"context", "dirtyInState", "oldState"},
	}
	self.Visit(parser, *parser.Entries[0], &current, template)

	self.renderers = append(self.renderers, self.CreateRenderer(current))

	js := template[parser.ScriptSource.StartIndex:parser.ScriptSource.EndIndex]

	for i, j := 0, len(self.renderers)-1; i < j; i, j = i+1, j-1 {
		self.renderers[i], self.renderers[j] = self.renderers[j], self.renderers[i]
	}

	code := `var component = function(options) {
		` + js + `; var currentComponent = null;` + strings.Join(self.renderers, "\n\n") +
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
				oldState = Object.assign({}, state)
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

	return code
}
