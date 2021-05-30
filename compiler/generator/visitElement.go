package generator

import (
	"elljo/compiler/parser"
	"strconv"
	"strings"
)

func (self *Generator) VisitElement(parser parser.Parser, children parser.Entry, current *Fragment) *Fragment {
	current.Counters.Element++
	name := "element_" + strconv.Itoa(current.Counters.Element)

	isComponent := false
	for _, componentImport := range parser.ScriptSource.Imports {
		if componentImport == children.Name {
			isComponent = true
		}
	}
	if isComponent {
		template := "$name$({target: $target$});"
		variables := map[string]string{
			"name":   children.Name,
			"target": current.Target,
		}
		initStatement := Statement{
			source:   self.BuildString(template, variables),
			mappings: [][]int{},
		}
		current.InitStatements = append(current.InitStatements, initStatement)
	} else {
		template := "var $name$ = document.createElement('$childrenName$');"
		variables := map[string]string{
			"name":         name,
			"childrenName": strings.ReplaceAll(children.Name, "\n", ""),
		}
		createStatement := self.BuildString(template, variables)
		mappings := [][]int{{}}

		for attributeIndex, attribute := range children.Attributes {
			if attribute.IsExpression {
				if attribute.IsCall {
					attributeCreateStatement := `var $name$_attr_$index$ = () => {
						return $value$
					};
					$name$.setAttribute("$attributeName$", $name$_attr_$index$());`

					mappings = append(mappings, []int{}, []int{}, []int{}, []int{}, []int{})

					variables := map[string]string{
						"name":          name,
						"index":         strconv.Itoa(attributeIndex),
						"value":         attribute.Value,
						"attributeName": attribute.Name,
					}

					createStatement += self.BuildString(attributeCreateStatement, variables)

					attributeUpdateStatementSource := `$name$.setAttribute("$attributeName$", $name$_attr_$index$());`

					attributeUpdateStatement := Statement{
						source:   self.BuildString(attributeUpdateStatementSource, variables),
						mappings: [][]int{},
					}

					current.UpdateStatments = append(current.UpdateStatments, attributeUpdateStatement)
				} else {
					variableCreateStatement := `$name$.setAttribute("$attributeName$", $value$);`
					variables := map[string]string{
						"name":          name,
						"attributeName": attribute.Name,
						"value":         attribute.Value,
					}
					mappings = append(mappings, []int{0, 0, children.Line, 0})
					createStatement += self.BuildString(variableCreateStatement, variables)
					variableUpdateStatementSource := `if(currentComponent.$value$IsDirty) {
								$name$.setAttribute("$attributeName$", $value$);
							}`

					variableUpdateStatement := Statement{
						source:   self.BuildString(variableUpdateStatementSource, variables),
						mappings: [][]int{{}, {0, 0, children.Line, 0}, {}},
					}

					current.UpdateStatments = append(current.UpdateStatments, variableUpdateStatement)
				}

			} else if attribute.IsEvent {
				if attribute.IsCall {
					attributeCreateStatement := `$name$.addEventListener("$attributeName$", () => {
						$value$
					});`
					mappings = append(mappings, []int{}, []int{}, []int{}, []int{})
					variables := map[string]string{
						"attributeName": attribute.Name,
						"name":          name,
						"value":         attribute.Value,
					}
					createStatement += self.BuildString(attributeCreateStatement, variables)
				} else {
					createStatement += name + `.addEventListener("` + attribute.Name + `", ` + attribute.Value + `);`
					mappings = append(mappings, []int{})
				}

			} else {
				if attribute.HasValue {
					createStatement += name + `.setAttribute("` + attribute.Name + `", "` + attribute.Value + `");`
					mappings = append(mappings, []int{})
				} else {
					createStatement += name + `.setAttribute("` + attribute.Name + `", "true");`
					mappings = append(mappings, []int{})
				}
			}
		}
		current.InitStatements = append(current.InitStatements, Statement{
			source:   createStatement,
			mappings: mappings,
		})

		removeStatementSource := name + ".parentNode.removeChild(" + name + ")"

		removeStatement := Statement{
			source:   removeStatementSource,
			mappings: [][]int{{}},
		}

		current.TeardownStatements = append(current.TeardownStatements, removeStatement)
	}

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
		IsComponent:        isComponent,
		UpdateContextChain: current.UpdateContextChain,
	}
}

func (self *Generator) VisitElementAfter(current *Fragment) {
	name := current.Target
	isComponent := current.IsComponent
	current.Parent.InitStatements = current.InitStatements
	current.Parent.UpdateStatments = current.UpdateStatments
	current.Parent.Counters = current.Counters
	current = current.Parent

	if !isComponent {
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
