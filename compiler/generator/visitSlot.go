package generator

import (
	"elljo/compiler/parser"
	"strconv"
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

	if name == "render" {
		current.Counters.Slot++
		name = "slot_" + strconv.Itoa(current.Counters.Slot)

		current.Slots = append(current.Slots, SlotEntry{
			Slot:     "default",
			Renderer: "render" + name,
		})
	} else {
		current.Parent.Slots = append(current.Parent.Slots, SlotEntry{
			Slot:     children.Parameter,
			Renderer: "render" + name,
		})
	}

	initStatement := "const render" + name + " = (target) => {"

	if len(current.SlotElements) > 0 {
		for _, entry := range current.SlotElements {
			initStatement += `
				target.appendChild(` + entry + `);`
		}

		initStatement += `
				return target;
			}`
	}

	current.InitStatements = append(current.InitStatements, Statement{
		source: initStatement,
		// TODO: Fix source mapping
		mappings: [][]int{{}},
	})

	current.Parent.InitStatements = append(current.Parent.InitStatements, current.InitStatements...)
}
