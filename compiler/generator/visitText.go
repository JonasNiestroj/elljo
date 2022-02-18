package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/utils"
	"strconv"
)

func (self *Generator) VisitText(children parser.Entry, current *Fragment) *Fragment {
	createStatementSource := `var text_$counter$ = document.createTextNode('$text$');
		$target$.appendChild(text_$counter$)`

	variables := map[string]string{
		"counter": strconv.Itoa(current.Counters.Text),
		"text":    children.Data,
		"target":  current.Target,
	}

	createStatement := Statement{
		source:   utils.BuildString(createStatementSource, variables),
		mappings: [][]int{{}},
	}

	if current.IsComponent {
		current.SlotElements = append(current.SlotElements, "text_"+strconv.Itoa(current.Counters.Text))
	}

	current.Counters.Text++

	current.InitStatements = append(current.InitStatements, createStatement)
	return current
}
