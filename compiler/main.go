package main

import (
	"elljo/compiler/generator"
	"elljo/compiler/parser"
	"elljo/compiler/service"
	"elljo/compiler/ssr"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	if os.Args[1] == "--service" {
		service.RunService()
		return
	}
	inputFile := os.Args[2]
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var parserInstance = parser.Parser{
		Index:    0,
		Template: string(data),
		Entries:  []*parser.Entry{},
	}

	parserInstance.Parse()

	if len(parserInstance.Errors) > 0 {
		for _, error := range parserInstance.Errors {
			println("Line " + strconv.Itoa(error.Line) + ": " + error.Message)
		}
		return
	}
	if os.Args[1] == "generate" {
		var generatorInstance = generator.Generator{
			FileName: "Component",
		}

		generated := generatorInstance.Generate(parserInstance, parserInstance.Template)

		indexFile := os.Args[3]
		index, err := ioutil.ReadFile(indexFile)
		if err != nil {
			panic(err)
		}
		output := strings.Replace(string(index), "{SCRIPT}", generated.Output, 1)

		if generated.Css != "" {
			cssTemplate := `
			<style>
				` + generated.Css + `
			</style>`
			output = strings.Replace(output, "{STYLE}", cssTemplate, 1)
		} else {
			output = strings.Replace(output, "{STYLE}", "", 1)
		}
		ioutil.WriteFile(os.Args[4], []byte(output), 0644)

		if len(os.Args) > 5 {
			ioutil.WriteFile(os.Args[5], []byte(generated.Sourcemap), 0644)
		}
	} else if os.Args[1] == "render" {
		var renderInstance = ssr.SSR{}

		ssrOutput := renderInstance.SSR(parserInstance)

		index, err := ioutil.ReadFile("./ssrindex.html")
		if err != nil {
			panic(err)
		}

		output := strings.Replace(string(index), "{HTML}", ssrOutput.Html, 1)

		if ssrOutput.Css != "" {
			cssTemplate := `
			<style>
				` + ssrOutput.Css + `
			</style>`
			output = strings.Replace(output, "{STYLE}", cssTemplate, 1)
		}

		ioutil.WriteFile(os.Args[3], []byte(output), 0644)
	}
}
