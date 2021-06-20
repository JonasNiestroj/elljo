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
	current := Fragment{
		UseAnchor:          false,
		Name:               "render",
		Target:             "target",
		InitStatements:     []Statement{},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		ContextChain:       []string{"context", "dirtyInState", "oldState"},
	}

	js := parser.ScriptSource.Source

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

	variables := ""
	setIsDirtyFalse := ""

	for _, variable := range parser.ScriptSource.Variables {
		propertyUpdate := ""
		for _, componentProperties := range self.componentProperties {
			if id, ok := componentProperties.Properties[variable]; ok {
				propertyUpdate += `
					this['component-` + strconv.Itoa(componentProperties.Index) + `'].$props['` + id + `'] = value`
			}
		}
		variables += `
			Object.defineProperty(this, "` + variable + `", {
				get() {
					return ` + variable + `;
				},
				set(value) {
					currentComponent.oldState["` + variable + `"] = ` + variable + `;
					` + variable + ` = value;
					currentComponent.` + variable + `IsDirty = true;
					new Observer(value, "` + variable + `")
					currentComponent.queueUpdate(); ` + propertyUpdate + `
				}
			})
		`
		setIsDirtyFalse += `
			currentComponent.` + variable + `IsDirty = false;`
	}

	properties := ""

	if len(parser.ScriptSource.Properties) > 0 {
		properties = "let that = this;"
	}

	for _, property := range parser.ScriptSource.Properties {
		properties += `
			Object.defineProperty(component.$props, "` + property + `", {
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

	code += `var component = function(options, props, events) {
		let currentComponent = null;` + js + variables +
		`; 
		` + strings.Join(renderersSources, "\n") +
		`
			let component = {};
			var state = {};
			component.oldState = {};
			var updating = false;
			var dirtyInState = [];

			component.$props = {};
			` + properties + `

			component.$events = {}
			if(events) {
				Object.keys(events).forEach(event => {
					if(!component.$events[event]) {
						component.$events[event] = [events[event]]
					} else {
						component.$events[event].push(events[event])
					}
				})
			}

			component.queueUpdate = function performUpdate() {
				if(!updating) {
					updating = true;
					Promise.resolve().then(component.update)
				}
			}
			component.update = function update() {
				updating = false;
				mainFragment.update();
				component.oldState = {}` + setIsDirtyFalse + `
			}
			component.teardown = function teardown () {
				mainFragment.teardown();
				mainFragment = null;
				state = {};
			};
			this.$event = function(name) {
				var callbacks = component.$events[name]
				if(callbacks) {
					const args = []
					for(let i = 1; i < arguments.length; i++) {
						args.push(arguments[i])
					}
					callbacks.forEach(callback => callback(...args))
				}
			}
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
			currentComponent = component;
			currentComponent.queueUpdate();

			function patchArray(array, name) {
				const methodsToPatch = ['push', 'pop', 'splice', 'sort', 'reverse', 'shift', 'unshift', 'fill']
				methodsToPatch.forEach(method => {
					const currentMethod = array[method]
					Object.defineProperty(array, method,  {
						enumerable: false,
						configurable: false,
						writable: false,
						value: function() {
							const result = currentMethod.apply(this, arguments)
							currentComponent[name + 'IsDirty'] = true;
							currentComponent.oldState[name] = array;
							currentComponent.queueUpdate();
							new Observer(result, name);
							return result;
						}
					});
				});
        	}

        	function Observer(value, name) {
			    if(!value || value.__observer__ ||(!Array.isArray(value) && typeof value !== 'object')) {
			        return
                }
                this.value = value
                value.__observer__ = this
                if(Array.isArray(value)) {
                    for(var i = 0; i < value.length; i++) {
                        new Observer(value[i], name)
                    }
                } else if (typeof value === 'object' && value !== null) {
                    const keys = Object.keys(value)
                    for(var i = 0; i < keys.length; i++) {
                        const key = keys[i]
                        // Check if current property is configurable
                        var prop = Object.getOwnPropertyDescriptor(value, key)
                        if((prop && !prop.configurable) || key === '__observer__') {
                            continue
                        }
                		let keyValue = value[key];
                        // TODO: Check for already existing getter/setter
                        new Observer(keyValue, name)
                        Object.defineProperty(value, key, {
                            enumerable: true,
                            configurable: true,
                            get: function() {
                                return keyValue
                            },
                            set: function(newValue) {
                                currentComponent[name + 'IsDirty'] = true;
                                currentComponent.oldState[name] = keyValue;
                                currentComponent.queueUpdate();
                                keyValue = newValue
                                new Observer(newValue, name)
                            }
                        })
                    }
                }
            }

        	const propertyNames = Object.getOwnPropertyNames(this)

			for(let i = 0; i < propertyNames.length; i++) {
				const property = propertyNames[i]
            	if(Array.isArray(this[property])) {
                	patchArray(this[property], property)
					new Observer(this[property], property)
            	} else {
            	    new Observer(this[property], property)
                }
			}

			return component;
	}
	export default component`

	style := parser.StyleSource.Source
	indexToAdd := 0
	for _, rule := range parser.StyleSource.Rules {
		style = style[0:rule.StartIndex+indexToAdd] + rule.Selector + "[scope-" + parser.StyleSource.Id + "]" + style[rule.EndIndex+indexToAdd:]
		indexToAdd += 8 + len(parser.StyleSource.Id)
	}

	return GeneratorOutput{Output: code, Sourcemap: sourcemap.CreateSourcemap(parser.FileName, mappingsStrings), Css: style}
}
