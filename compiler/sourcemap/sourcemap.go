package sourcemap

type Sourcemap struct {
	Mappings []string
}

func CreateSourcemap(values [][]int) string {
	// Add an empty line because the first line creates a function
	mapping := "\";"
	for _, value := range values {
		mapping += EncodeValues(value) + ";"
	}
	return `{"version": 3, "names": [], "mappings": ` + mapping + `", "sources": ["index.jo"]}`
}
