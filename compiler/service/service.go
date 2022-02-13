package service

import (
	"bufio"
	"elljo/compiler/generator"
	"elljo/compiler/parser"
	"encoding/json"
	"os"
	"strings"
)

func RunService() {
	scanner := bufio.NewScanner(os.Stdin)
	output := bufio.NewWriter(os.Stdout)

	message := ""

	for scanner.Scan() {
		message += scanner.Text() + "\n"
	}

	if strings.HasPrefix(message, "compile") {
		fileName := strings.Split(message, " ")[1]
		var parserInstance = parser.Parser{
			Index:    0,
			Template: strings.Replace(message, "compile "+fileName+" ", "", 1),
			Entries:  []*parser.Entry{},
			FileName: fileName,
		}
		parserInstance.Parse()

		if len(parserInstance.Errors) > 0 {
			errors, err := json.Marshal(parserInstance.Errors)
			if err != nil {
				output.Write([]byte(err.Error()))
				output.Flush()
				panic(err)
			}
			output.Write(errors)
			output.Flush()
			return
		}

		generatorInstance := generator.Generator{
			FileName: strings.Split(fileName, ".")[0],
		}

		generated := generatorInstance.Generate(parserInstance, parserInstance.Template)

		outputJson, err := json.Marshal(generated)

		if err != nil {
			output.Write([]byte(err.Error()))
			output.Flush()
			panic(err)
		}
		output.Write(outputJson)
		output.Flush()

	}
}
