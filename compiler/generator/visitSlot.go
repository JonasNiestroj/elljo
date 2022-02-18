package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/utils"
	"strconv"
	"strings"
)

func (self *Generator) VisitSlot(children parser.Entry, current *Fragment) *Fragment {
	current.Counters.Slot++
	name := "slot_" + strconv.Itoa(current.Counters.Slot)

	return &Fragment{
		UseAnchor:          true,
		Name:               name,
		Target:             "target",
		ContextChain:       current.ContextChain,
		InitStatements:     []Statement{},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		Counters:           current.Counters,
		Parent:             current,
		UpdateContextChain: current.UpdateContextChain,
		IsComponent:        true,
	}
}

func (self *Generator) VisitSlotAfter(children parser.Entry, current *Fragment) {

	name := current.Name

	if strings.HasPrefix(name, "render") {
		current.Counters.Slot++
		name = "slot_" + strconv.Itoa(current.Counters.Slot)

		current.Slots = append(current.Slots, SlotEntry{
			Slot:     "'default'",
			Renderer: "render" + name,
		})
	} else {
		current.Parent.Slots = append(current.Parent.Slots, SlotEntry{
			Slot:     children.Parameter,
			Renderer: "render" + name,
		})
	}

	initStatement := `const render$name$ = () => {
		return {
			render: (target) => {
				$slots$
			},
			teardown: () => {
				$teardown$
			}
		}
	}`

	slots := ""
	teardown := ""

	if len(current.SlotElements) > 0 {
		for _, entry := range current.SlotElements {
			slots += `
				target.appendChild(` + entry + `);`

			teardown += `
				` + entry + `.parentNode.removeChild(` + entry + `);`
		}

		slots += `
				return target;`
	}

	if current.IsComponent && !strings.HasPrefix(name, "render") {
		if !utils.IsOnlyStringExpression(children.Expression) && children.Parameter != "" {
			variableUpdateStatementSource := `console.log(this.$variablesToUpdate);if(this.$variablesToUpdate.includes('$value$')) {
								this['component-$componentIndex$'].$updateSlot($value$, this.oldState["$value$"]); 
							}`

			variables := map[string]string{
				"value":          children.Parameter,
				"componentIndex": strconv.Itoa(self.componentCounter),
			}

			variableUpdateStatement := Statement{
				source:   utils.BuildString(variableUpdateStatementSource, variables),
				mappings: [][]int{{}, {0, 0, children.Line, 0}, {}},
			}

			current.UpdateStatments = append(current.UpdateStatments, variableUpdateStatement)
		}
	}

	variables := map[string]string{
		"slots":    slots,
		"teardown": teardown,
		"name":     name,
	}

	current.InitStatements = append(current.InitStatements, Statement{
		source: utils.BuildString(initStatement, variables),
		// TODO: Fix source mapping
		mappings: [][]int{{}},
	})

	current.Parent.InitStatements = append(current.Parent.InitStatements, current.InitStatements...)
	current.Parent.UpdateStatments = append(current.Parent.UpdateStatments, current.UpdateStatments...)
}
