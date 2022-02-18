package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/utils"
	"strconv"
	"strings"
	"unicode"
)

type ComponentProperties struct {
	Index      int
	Properties map[string]string
}

func (self *Generator) handleComponent(children parser.Entry, current *Fragment, isGlobalComponent bool, name string) {
	componentProperties := ComponentProperties{
		Index:      self.componentCounter,
		Properties: map[string]string{},
	}
	self.componentCounter++

	props := ""
	events := ""

	createTemplate := ""

	for _, attribute := range children.Attributes {
		if !attribute.IsEvent {
			if props != "" {
				props += ", '" + attribute.Name + "': " + attribute.Value
			} else {
				props += attribute.Name + ": " + attribute.Value
			}
			createTemplate += `
					if(!this.$propsBindings["` + attribute.Name + `"]) {
						this.$propsBindings["` + attribute.Name + `"] = [];							
					}
					this.$propsBindings["` + attribute.Name + `"].push('component-` + strconv.Itoa(componentProperties.Index) + `');
`
			componentProperties.Properties[attribute.Value] = attribute.Name
		} else {
			if attribute.IsCall {
				events += "'" + attribute.Name + "': () => {" + attribute.Value + "},"
			} else {
				events += "'" + attribute.Name + "': " + attribute.Value + ","
			}
		}
	}

	self.componentProperties = append(self.componentProperties, componentProperties)

	if isGlobalComponent {
		createTemplate += `let $element_name$ = window.__elljo__.components["$name$"]
								if(!$element_name$) {
									//TODO: Log error and better error handling (only return this component)
									return
								}
`
	}

	createTemplate += `this['component-` + strconv.Itoa(componentProperties.Index) +
		`'] = new $element_name$({target: $target$`

	if len(current.SlotElements) > 0 || len(current.Slots) > 0 {
		createTemplate += ", slots"
	}

	createTemplate += `}, {` + props + `}, {` + events + `});`
	variables := map[string]string{
		"name":         children.Name,
		"target":       current.Parent.Target,
		"props":        props,
		"element_name": name,
	}
	initStatement := Statement{
		source:   utils.BuildString(createTemplate, variables),
		mappings: [][]int{},
	}
	current.InitStatements = append(current.InitStatements, initStatement)

	removeStatementSource := "this['component-" + strconv.Itoa(componentProperties.Index) + "'].teardown();"

	removeStatement := Statement{
		source:   removeStatementSource,
		mappings: [][]int{{}},
	}

	current.TeardownStatements = append(current.TeardownStatements, removeStatement)
}

func (self *Generator) VisitElement(parser parser.Parser, children parser.Entry, current *Fragment) *Fragment {

	current.Counters.Element++
	name := "element_" + strconv.Itoa(current.Counters.Element)

	if current.IsComponent {
		current.SlotElements = append(current.SlotElements, name)
	}

	isComponent := false
	isGlobalComponent := false
	for _, componentImport := range parser.ScriptSource.Imports {
		if componentImport.Name == children.Name {
			isComponent = true
			name = componentImport.Name
		}
	}

	if !isComponent {
		if unicode.IsUpper(rune(children.Name[0])) {
			isGlobalComponent = true
		}
	}

	if !isComponent && !isGlobalComponent {
		var template strings.Builder

		if len(children.LoopIndices) > 0 {

			// Indicates whether the element has an expression, a call or an event as an attribute
			hasRelevantAttribute := false
			for _, attribute := range children.Attributes {
				if attribute.IsExpression || attribute.IsEvent {
					hasRelevantAttribute = true
				}
			}

			// If all attributes are static we can return everything
			if !hasRelevantAttribute {
				return &Fragment{
					Target:             name,
					TeardownStatements: current.TeardownStatements,
					Name:               current.Name,
					InitStatements:     current.InitStatements,
					Counters:           current.Counters,
					ContextChain:       current.ContextChain,
					UpdateStatments:    current.UpdateStatments,
					UseAnchor:          current.UseAnchor,
					Parent:             current,
					IsComponent:        isComponent || isGlobalComponent,
					UpdateContextChain: current.UpdateContextChain,
				}
			}

			template.WriteString("\nvar $name$ = $initHtml")

			for _, index := range children.LoopIndices {
				template.WriteString(".childNodes[" + strconv.Itoa(index) + "]")
			}

			template.WriteString(";\n")
		} else {
			if children.Namespace == "" {
				template.WriteString("\nvar $name$ = elementCache.$childrenName$.cloneNode(true);\n")
			} else {
				template.WriteString("\nvar $name$ = document.createElementNS(\"$childrenNamespace$\", \"$childrenName$\");\n")
			}
		}

		childrenName := strings.ReplaceAll(children.Name, "\n", "")
		self.elements[childrenName] = struct{}{}

		variables := map[string]string{
			"name":              name,
			"childrenName":      childrenName,
			"childrenNamespace": children.Namespace,
		}

		createStatement := utils.BuildString(template.String(), variables)
		mappings := [][]int{{}}

		for _, attribute := range children.Attributes {
			if attribute.IsExpression {
				variableCreateStatement := `$name$.setAttribute("$attributeName$", $value$);`
				variables := map[string]string{
					"name":          name,
					"attributeName": attribute.Name,
					"value":         attribute.Value,
				}
				mappings = append(mappings, []int{0, 0, children.Line, 0})
				createStatement += utils.BuildString(variableCreateStatement, variables)
				variableUpdateStatementSource := `if(this.$variablesToUpdate.includes($value$)) {
								$name$.setAttribute("$attributeName$", $value$);
							}`

				variableUpdateStatement := Statement{
					source:   utils.BuildString(variableUpdateStatementSource, variables),
					mappings: [][]int{{}, {0, 0, children.Line, 0}, {}},
				}

				current.UpdateStatments = append(current.UpdateStatments, variableUpdateStatement)

			} else if attribute.IsEvent {
				attributeCreateStatement := ""

				if attribute.IsCall {
					attributeCreateStatement = `$name$.addEventListener("$attributeName$", () => {
						$value$
					});`
				} else {
					attributeCreateStatement = `$name$.addEventListener("$attributeName$", $value$)`
				}

				mappings = append(mappings, []int{}, []int{}, []int{}, []int{})
				variables := map[string]string{
					"attributeName": attribute.Name,
					"name":          name,
					"value":         attribute.Value,
				}
				createStatement += utils.BuildString(attributeCreateStatement, variables)

			} else {
				if attribute.HasValue {
					createStatement += name + `.setAttribute("` + attribute.Name + `", ` + attribute.Value + `);`
					mappings = append(mappings, []int{})
				} else {
					createStatement += name + `.setAttribute("` + attribute.Name + `", true);`
					mappings = append(mappings, []int{})
				}
			}
		}
		current.InitStatements = append(current.InitStatements, Statement{
			source:   createStatement,
			mappings: mappings,
		})

		if parser.StyleSource.Id != "" {
			scopeStatement := `$name$.setAttribute("$id$", "");`
			mappings = append(mappings, []int{})
			variables := map[string]string{
				"id":   "scope-" + parser.StyleSource.Id,
				"name": name,
			}
			current.InitStatements = append(current.InitStatements, Statement{
				source:   utils.BuildString(scopeStatement, variables),
				mappings: mappings,
			})
		}

		removeStatementSource := name + ".parentNode.removeChild(" + name + ")"

		removeStatement := Statement{
			source:   removeStatementSource,
			mappings: [][]int{{}},
		}

		current.TeardownStatements = append(current.TeardownStatements, removeStatement)
	}

	needsAppend := !(len(children.LoopIndices) > 0)

	return &Fragment{
		Target:             name,
		TeardownStatements: current.TeardownStatements,
		Name:               current.Name,
		InitStatements:     current.InitStatements,
		Counters:           current.Counters,
		ContextChain:       current.ContextChain,
		UpdateStatments:    current.UpdateStatments,
		UseAnchor:          current.UseAnchor,
		Parent:             current,
		IsComponent:        isComponent || isGlobalComponent,
		UpdateContextChain: current.UpdateContextChain,
		NeedsAppend:        needsAppend,
		SlotElements:       []string{},
	}
}

func (self *Generator) VisitElementAfter(parser parser.Parser, current *Fragment, children parser.Entry) {
	name := current.Target
	needsAppend := current.NeedsAppend
	isComponent := current.IsComponent

	if len(current.SlotElements) > 0 {
		self.VisitSlotAfter(children, current)
	}

	if len(current.Slots) > 0 {
		initStatementSource := "var slots = {"

		for _, slot := range current.Slots {
			initStatementSource += "[" + slot.Slot + "]: " + slot.Renderer + ","
		}

		initStatementSource += "};"

		initStatement := Statement{
			source: initStatementSource,
			// TODO: Fix source map
			mappings: [][]int{{}},
		}
		current.InitStatements = append(current.InitStatements, initStatement)
	}

	if isComponent {
		isGlobalComponent := true
		for _, componentImport := range parser.ScriptSource.Imports {
			if componentImport.Name == children.Name {
				isGlobalComponent = false
			}
		}

		self.handleComponent(children, current, isGlobalComponent, name)
	}

	current.Parent.InitStatements = current.InitStatements
	current.Parent.UpdateStatments = current.UpdateStatments
	current.Parent.Counters = current.Counters
	current = current.Parent

	if !isComponent && needsAppend && !current.IsComponent {
		if current.UseAnchor && current.Target == "target" {
			initStatement := Statement{
				source:   "target.insertBefore(" + name + ", anchor);",
				mappings: [][]int{{}},
			}
			current.InitStatements = append(current.InitStatements, initStatement)
		} else {
			initStatement := Statement{
				source:   current.Target + ".appendChild(" + name + ");",
				mappings: [][]int{{}},
			}
			current.InitStatements = append(current.InitStatements, initStatement)
		}
	}
}
