package sourcemap

import "strings"

type Sourcemap struct {
	Mappings []string
}

func CreateSourcemap(fileName string, mappings []string) string {
	return `{"version": 3, "names": [], "mappings": "` + strings.Join(mappings, "") + `", "sources": ["` + fileName + `"]}`
}
