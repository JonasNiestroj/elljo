package sourcemap

var (
	chars = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '/', '='}
)

func EncodeValues(values []int) string {
	encodedValue := ""
	for _, value := range values {
		encodedValue += EncodeValue(value)
	}

	return encodedValue
}

func EncodeValue(value int) string {
	encodedValue := ""

	if value < 0 {
		value = (-value << 1) | 1
	} else {
		value <<= 1
	}

	for {
		valueAnd := value & 31

		value >>= 5

		if value > 0 {
			valueAnd |= 32
		}

		encodedValue += string(chars[valueAnd])
		if value <= 0 {
			break
		}
	}

	return encodedValue
}
