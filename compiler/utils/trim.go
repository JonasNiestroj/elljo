package utils

func TrimStart(str string) string {
	index := 0
	currentChar := string(str[index])

	for currentChar == " " {
		index++
		currentChar = string(str[index])
	}

	return str[index:]
}

func TrimEnd(str string) string {
	index := len(str)
	currentChar := string(str[index])

	for currentChar == " " {
		index--
		currentChar = string(str[index])
	}

	return str[:index]
}
