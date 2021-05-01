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
	for scanner.Scan() {
		bytes := scanner.Bytes()
		message := strings.ReplaceAll(string(bytes), "\\n", "\n")
		if strings.HasPrefix(message, "compile") {
			var parserInstance = parser.Parser{
				Index:    0,
				Template: strings.Replace(message, "compile ", "", 1),
				Entries:  []*parser.Entry{},
			}
			parserInstance.Parse()

			generatorInstance := generator.Generator{}

			generated := generatorInstance.Generate(parserInstance, parserInstance.Template)

			outputJson, err := json.Marshal(generated)

			if err != nil {
				os.Stdout.Write([]byte(err.Error()))
				panic(err)
			}

			os.Stdout.Write(outputJson)

			break
		}
	}
}
