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
		current.InitStatements = append(current.InitStatements, self.BuildString(template, variables))
	} else {
		template := "var $name$ = document.createElement('$childrenName$');"
		variables := map[string]string{
			"name":         name,
			"childrenName": strings.ReplaceAll(children.Name, "\n", ""),
		}
		createStatement := self.BuildString(template, variables)

		for attributeIndex, attribute := range children.Attributes {
			if attribute.IsExpression {
				if attribute.IsCall {
					attributeCreateStatement := `var $name$_attr_$index$ = () => {
						$context$
						return $value$
					};
					$name$.setAttribute("$attributeName$", $name$_attr_$index$());`

					contextString := ""
					for _, context := range current.ContextChain {
						if context != "context" && context != "dirtyInState" && context != "oldState" {
							contextString += `var ` + context + ` = currentContext.` + context + `;`
						}
					}
					variables := map[string]string{
						"name":          name,
						"index":         strconv.Itoa(attributeIndex),
						"context":       contextString,
						"value":         attribute.Value,
						"attributeName": attribute.Name,
					}

					createStatement += self.BuildString(attributeCreateStatement, variables)

					attributeUpdateStatement := `$name$.setAttribute("$attributeName$", $name$_attr_$index$());`

					current.UpdateStatments = append(current.UpdateStatments, self.BuildString(attributeUpdateStatement, variables))
				} else {
					for _, variable := range parser.ScriptSource.Variables {
						if variable == attribute.Value {
							variableCreateStatement := `$name$.setAttribute("$attributeName$", $value$);`
							variables := map[string]string{
								"name":          name,
								"attributeName": attribute.Name,
								"value":         attribute.Value,
							}
							createStatement += self.BuildString(variableCreateStatement, variables)
							variableUpdateStatement := `if(dirtyInState.includes("$value")) {
								$name$.setAttribute("$attributeName$", context.$value$);
							}`
							current.UpdateStatments = append(current.UpdateStatments, self.BuildString(variableUpdateStatement, variables))
						}
					}
				}

			} else if attribute.IsEvent {
				if attribute.IsCall {
					attributeCreateStatement := `$name$.addEventListener("$attributeName$", () => {
						$context$
						$value$
					});`

					contextTemplate := ""
					for _, context := range current.ContextChain {
						if context != "context" && context != "dirtyInState" && context != "oldState" {
							contextTemplate += `var ` + context + ` = currentContext.` + context + `;\n`
						}
					}
					variables := map[string]string{
						"context":       contextTemplate,
						"attributeName": attribute.Name,
						"name":          name,
					}
					createStatement += self.BuildString(attributeCreateStatement, variables)
				} else {
					createStatement += name + `.addEventListener("` + attribute.Name + `", ` + attribute.Value + `);`
				}

			} else {
				if attribute.HasValue {
					createStatement += name + `.setAttribute("` + attribute.Name + `", "` + attribute.Value + `");`
				} else {
					createStatement += name + `.setAttribute("` + attribute.Name + `", "true");`
				}
			}
		}
		current.InitStatements = append(current.InitStatements, createStatement)

		removeStatement := name + ".parentNode.removeChild(" + name + ")"

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
			current.InitStatements = append(current.InitStatements, "target.insertBefore("+name+", anchor);")
		} else {
			current.InitStatements = append(current.InitStatements, current.Target+".appendChild("+name+");")
		}
	}
}
