package sourcemap

import "strings"

type Sourcemap struct {
	Mappings []string
}

func CreateSourcemap(mappings []string) string {
	return `{"version": 3, "names": [], "mappings": "` + strings.Join(mappings, "") + `", "sources": ["index.jo"]}`
}
