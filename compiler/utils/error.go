package utils

type Error struct {
	Line        int    `json:"line"`
	Message     string `json:"message"`
	StartColumn int    `json:"startColumn"`
	EndColumn   int    `json:"endColumn"`
}
