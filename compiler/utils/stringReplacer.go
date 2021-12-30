package utils

import "sort"

type Chunk struct {
	start int
	end   int
	text  string
}

type StringReplacer struct {
	Text   string
	chunks []Chunk
}

func (stringReplacer *StringReplacer) Replace(start int, end int, text string) {
	stringReplacer.chunks = append(stringReplacer.chunks, Chunk{
		start: start,
		end:   end,
		text:  text,
	})
}

func (stringReplacer *StringReplacer) String() string {
	sort.Slice(stringReplacer.chunks, func(i, j int) bool {
		return stringReplacer.chunks[i].start < stringReplacer.chunks[j].start
	})

	buildedString := ""
	lastEndIndex := 0

	if len(stringReplacer.chunks) == 0 {
		buildedString = stringReplacer.Text
		lastEndIndex = len(stringReplacer.Text)
	}

	for _, chunk := range stringReplacer.chunks {
		buildedString += stringReplacer.Text[lastEndIndex:chunk.start]
		buildedString += chunk.text
		lastEndIndex = chunk.end
	}

	if lastEndIndex != len(stringReplacer.Text) {
		buildedString += stringReplacer.Text[lastEndIndex:]
	}

	return buildedString
}
