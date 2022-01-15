package sourcemap

type Sourcemap struct {
	Mappings []string
}

func CreateSourcemap(fileName string, mappings string) string {
	return `{"version": 3, "names": [], "mappings": "` + mappings + `", "sources": ["` + fileName + `"]}`
}
