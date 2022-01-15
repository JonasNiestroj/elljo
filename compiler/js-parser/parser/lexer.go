package parser

import (
	"elljo/compiler/js-parser/token"
	"elljo/compiler/js-parser/unistring"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func (self *Parser) ScanIdentifier() (string, unistring.String, bool, error) {
	offset := self.CharOffset
	hasEscape := false
	isUnicode := false
	length := 0
	for IsIdentifierPart(self.Char) {
		r := self.Char
		length++
		if r == '\\' {
			hasEscape = true
			distance := self.CharOffset - offset
			self.Read()
			if self.Char != 'u' {
				return "", "", false, fmt.Errorf("Invalid identifer escape character: %c (%s)", self.Char, string(self.Char))
			}
			var value rune
			for j := 0; j < 4; j++ {
				self.Read()
				decimal, ok := Hex2Decimal(byte(self.Char))
				if !ok {
					return "", "", false, fmt.Errorf("Invalid identifier escape character: %c (%s)", self.Char, string(self.Char))
				}
				value = value<<4 | decimal
			}
			if value == '\\' {
				return "", "", false, fmt.Errorf("INvalid identifier escape value: %c (%s)", value, string(value))
			} else if distance == 0 {
				if !IsIdentifierStart(value) {
					return "", "", false, fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
				}
			} else if distance > 0 {
				if !IsIdentifierPart(value) {
					return "", "", false, fmt.Errorf("Invalid identifer escape value: %c (%s)", value, string(value))
				}
			}
			r = value
		}
		if r >= utf8.RuneSelf {
			isUnicode = true
			if r > 0xFFFF {
				length++
			}
		}
		self.Read()
	}

	literal := self.Template[offset:self.CharOffset]
	var parsed unistring.String
	if hasEscape || isUnicode {
		var err error
		parsed, err = ParseStringLiteral1(literal, length, isUnicode)
		if err != nil {
			return "", "", false, err
		}
	} else {
		parsed = unistring.String(literal)
	}
	return literal, parsed, hasEscape, nil
}

func IsLineWhiteSpace(chr rune) bool {
	switch chr {
	case '\u0009', '\u000b', '\u000c', '\u0020', '\u00a0', '\ufeff':
		return true
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return false
	case '\u0085':
		return false
	}

	return unicode.IsSpace(chr)
}

func IsLineTerminator(chr rune) bool {
	switch chr {
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}
	return false
}

func IsId(tkn token.Token) bool {
	switch tkn {
	case token.KEYWORD,
		token.BOOLEAN,
		token.NULL,
		token.THIS,
		token.IF,
		token.IN,
		token.OF,
		token.DO,
		token.VAR,
		token.FOR,
		token.NEW,
		token.TRY,
		token.ELSE,
		token.CASE,
		token.VOID,
		token.WITH,
		token.WHILE,
		token.BREAK,
		token.CATCH,
		token.THROW,
		token.RETURN,
		token.TYPEOF,
		token.DELETE,
		token.SWITCH,
		token.DEFAULT,
		token.FINALLY,
		token.FUNCTION,
		token.CONTINUE,
		token.DEBUGGER,
		token.INSTANCEOF:
		return true
	}
	return false
}

func (self *Parser) Scan() (tkn token.Token, literal string, parsedLiteral unistring.String, index int) {
	self.ImplicitSemicolon = false
	for {
		self.SkipWhiteSpace()
		index = self.IndexOf(self.CharOffset)
		insertSemicolon := false
		switch chr := self.Char; {
		case IsIdentifierStart(chr):
			var err error
			var hasEscape bool
			literal, parsedLiteral, hasEscape, err = self.ScanIdentifier()
			if err != nil {
				tkn = token.ILLEGAL
				break
			}
			if len(parsedLiteral) > 1 {
				tkn = token.StringIsKeyword(string(parsedLiteral))

				switch tkn {
				case 0:
					if parsedLiteral == "true" || parsedLiteral == "false" {
						if hasEscape {
							tkn = token.ILLEGAL
							return
						}
						self.InsertSemicolon = true
						tkn = token.BOOLEAN
						return
					} else if parsedLiteral == "null" {
						if hasEscape {
							tkn = token.ILLEGAL
							return
						}
						self.InsertSemicolon = true
						tkn = token.NULL
						return
					}
				case token.KEYWORD:
					if hasEscape {
						tkn = token.ILLEGAL
						return
					}
					tkn = token.KEYWORD
					return
				case token.THIS,
					token.BREAK,
					token.THROW,
					token.RETURN,
					token.CONTINUE,
					token.DEBUGGER:
					if hasEscape {
						tkn = token.ILLEGAL
						return
					}
					self.InsertSemicolon = true
					return
				default:
					if hasEscape {
						tkn = token.ILLEGAL
					}
					return
				}
			}
			self.InsertSemicolon = true
			tkn = token.IDENTIFIER
			return
		case '0' <= chr && chr <= '9':
			self.InsertSemicolon = true
			tkn, literal = self.ScanNumericLiteral(false)
			return
		default:
			self.Read()
			switch chr {
			case -1:
				if self.InsertSemicolon {
					self.InsertSemicolon = false
					self.ImplicitSemicolon = true
				}
				tkn = token.EOF
			case '\r', '\n', '\u2028', '\u2029':
				self.InsertSemicolon = false
				self.ImplicitSemicolon = true
				continue
			case ':':
				tkn = token.COLON
			case '.':
				if DigitValue(self.Char) < 10 {
					insertSemicolon = true
					tkn, literal = self.ScanNumericLiteral(true)
				} else {
					if self.Char == '.' {
						self.Read()
						if self.Char == '.' {
							self.Read()
							tkn = token.Ellipsis
						} else {
							tkn = token.PERIOD
						}
					} else {
						tkn = token.PERIOD
					}
				}
			case ',':
				tkn = token.COMMA
			case ';':
				tkn = token.SEMICOLON
			case '(':
				tkn = token.LEFT_PARENTHESIS
			case ')':
				tkn = token.RIGHT_PARENTHESIS
				insertSemicolon = true
			case '[':
				tkn = token.LEFT_BRACKET
			case ']':
				tkn = token.RIGHT_BRACKET
				insertSemicolon = true
			case '{':
				tkn = token.LEFT_BRACE
			case '}':
				tkn = token.RIGHT_BRACE
				insertSemicolon = true
			case '+':
				tkn = self.Switch3(token.PLUS, token.ADD_ASSIGN, '+', token.INCREMENT)
				if tkn == token.INCREMENT {
					insertSemicolon = true
				}
			case '-':
				tkn = self.Switch3(token.MINUS, token.SUBTRACT_ASSIGN, '-', token.DECREMENT)
				if tkn == token.DECREMENT {
					insertSemicolon = true
				}
			case '*':
				tkn = self.Switch2(token.MULTIPLY, token.MULTIPLY_ASSIGN)
			case '/':
				if self.Char == '/' {
					self.SkipSingleLineComment()
					continue
				} else if self.Char == '*' {
					self.SkipMultiLineComment()
					continue
				} else {
					tkn = self.Switch2(token.SLASH, token.QUOTIENT_ASSIGN)
					insertSemicolon = true
				}
			case '%':
				tkn = self.Switch2(token.REMAINDER, token.REMAINDER_ASSIGN)
			case '^':
				tkn = self.Switch2(token.EXCLUSIVE_OR, token.EXCLUSIVE_OR_ASSIGN)
			case '<':
				tkn = self.Switch4(token.LESS, token.LESS_OR_EQUAL, '<', token.SHIFT_LEFT, token.SHIFT_LEFT_ASSIGN)
			case '>':
				tkn = self.Switch6(token.GREATER, token.GREATER_OR_EQUAL, '>', token.SHIFT_RIGHT, token.SHIFT_RIGHT_ASSIGN, '>', token.UNSIGNED_SHIFT_RIGHT, token.UNSIGNED_SHIFT_RIGHT_ASSIGN)
			case '=':
				tkn = self.Switch3(token.ASSIGN, token.EQUAL, '>', token.ARROW_FUNCTION)
				if tkn == token.EQUAL && self.Char == '=' {
					self.Read()
					tkn = token.STRICT_EQUAL
				}
			case '!':
				tkn = self.Switch2(token.NOT, token.NOT_EQUAL)
				if tkn == token.NOT_EQUAL && self.Char == '=' {
					self.Read()
					tkn = token.STRICT_NOT_EQUAL
				}
			case '&':
				tkn = self.Switch3(token.AND, token.AND_ASSIGN, '&', token.LOGICAL_AND)
			case '|':
				tkn = self.Switch3(token.OR, token.OR_ASSIGN, '|', token.LOGICAL_OR)
			case '~':
				tkn = token.BITWISE_NOT
			case '?':
				tkn = token.QUESTION_MARK
			case '"', '\'':
				insertSemicolon = true
				tkn = token.STRING
				var err error
				literal, parsedLiteral, err = self.ScanString(self.CharOffset-1, true)
				if err != nil {
					tkn = token.ILLEGAL
				}
			default:
				self.ErrorUnexpected(chr, index)
				tkn = token.ILLEGAL
			}
		}
		self.InsertSemicolon = insertSemicolon
		return
	}
}

func (self *Parser) Switch2(tkn0, tkn1 token.Token) token.Token {
	if self.Char == '=' {
		self.Read()
		return tkn1
	}
	return tkn0
}

func (self *Parser) Switch3(tkn0, tkn1 token.Token, chr2 rune, tkn2 token.Token) token.Token {
	if self.Char == '=' {
		self.Read()
		return tkn1
	}
	if self.Char == chr2 {
		self.Read()
		return tkn2
	}
	return tkn0
}

func (self *Parser) Switch4(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token) token.Token {
	if self.Char == '=' {
		self.Read()
		return tkn1
	}
	if self.Char == chr2 {
		self.Read()
		if self.Char == '=' {
			self.Read()
			return tkn3
		}
		return tkn2
	}
	return tkn0
}

func (self *Parser) Switch6(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token, chr3 rune, tkn4, tkn5 token.Token) token.Token {
	if self.Char == '=' {
		self.Read()
		return tkn1
	}
	if self.Char == chr2 {
		self.Read()
		if self.Char == '=' {
			self.Read()
			return tkn3
		}
		if self.Char == chr3 {
			self.Read()
			if self.Char == '=' {
				self.Read()
				return tkn5
			}
			return tkn4
		}
		return tkn2
	}
	return tkn0
}

func (self *Parser) Peek() rune {
	if self.Offset+1 < self.Length {
		return rune(self.Template[self.Offset+1])
	}
	return -1
}
func (self *Parser) Read() {
	if self.Offset < self.Length {
		self.CharOffset = self.Offset
		chr, width := rune(self.Template[self.Offset]), 1
		if chr >= utf8.RuneSelf {
			chr, width = utf8.DecodeRuneInString(self.Template[self.Offset:])
			if chr == utf8.RuneError && width == 1 {
				self.Error(self.Index, "Invalid UTF-8 character")
			}
		}
		self.Offset += width
		self.Char = chr
	} else {
		self.CharOffset = self.Length
		self.Char = -1
	}
}

func (self *Parser) SkipSingleLineComment() {
	for self.Char != -1 {
		self.Read()
		if IsLineTerminator(self.Char) {
			return
		}
	}
}

func (self *Parser) SkipMultiLineComment() {
	start := self.Index
	self.Read()
	for self.Char >= 0 {
		chr := self.Char
		self.Read()
		if chr == '*' && self.Char == '/' {
			self.Read()
			return
		}
	}
	self.ErrorUnexpected(self.Char, start)
}

func (self *Parser) SkipWhiteSpace() {
	for {
		switch self.Char {
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			self.Read()
			continue
		case '\r':
			if self.Peek() == '\n' {
				self.Read()
			}
			fallthrough
		case '\u2028', '\u2029', '\n':
			if self.InsertSemicolon {
				return
			}
			self.Read()
			continue
		}
		if self.Char >= utf8.RuneSelf {
			if unicode.IsSpace(self.Char) {
				self.Read()
				continue
			}
		}
		break
	}
}

func (self *Parser) SkipLineWhiteSpace() {
	for IsLineWhiteSpace(self.Char) {
		self.Read()
	}
}

func (self *Parser) ScanMantissa(base int) {
	for DigitValue(self.Char) < base {
		self.Read()
	}
}

func (self *Parser) ScanEscape(quote rune) (int, bool) {
	var length, base uint32
	chr := self.Char
	switch chr {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		length, base = 3, 8
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'':
		self.Read()
		return 1, false
	case '\r':
		self.Read()
		if self.Char == '\n' {
			self.Read()
			return 2, false
		}
		return 1, false
	case '\n':
		self.Read()
		return 1, false
	case '\u2028', '\u2029':
		self.Read()
		return 1, true
	case 'x':
		self.Read()
		length, base = 2, 16
	case 'u':
		self.Read()
		length, base = 4, 16
	default:
		self.Read()
	}

	if length > 0 {
		var value uint32
		for ; length > 0 && self.Char != quote && self.Char >= 0; length-- {
			digit := uint32(DigitValue(self.Char))
			if digit >= base {
				break
			}
			value = value*base + digit
			self.Read()
		}
		chr = rune(value)
	}
	if chr >= utf8.RuneSelf {
		if chr > 0xFFFF {
			return 2, true
		}
		return 1, true
	}
	return 1, false
}

func (self *Parser) ScanString(offset int, parse bool) (literal string, parsed unistring.String, err error) {
	quote := rune(self.Template[offset])
	length := 0
	isUnicode := false
	for self.Char != quote {
		chr := self.Char
		if chr == '\n' || chr == '\r' || chr == '\u2028' || chr == '\u2029' || chr < 0 {
			return self.CheckNewLine(quote)
		}
		self.Read()
		if chr == '\\' {
			if self.Char == '\n' || self.Char == '\r' || self.Char == '\u2028' || self.Char == '\u2029' || self.Char < 0 {
				if quote == '/' {
					return self.CheckNewLine(quote)
				}
				self.ScanNewline()
			} else {
				l, u := self.ScanEscape(quote)
				length += l
				if u {
					isUnicode = true
				}
			}
			continue
		} else if chr == '[' && quote == '/' {
			quote = -1
		} else if chr == ']' && quote == -1 {
			quote = '/'
		}
		if chr >= utf8.RuneSelf {
			isUnicode = true
			if chr > 0xFFFF {
				length++
			}
		}
		length++
	}
	self.Read()
	literal = self.Template[offset:self.CharOffset]
	if parse {
		parsed, err = ParseStringLiteral1(literal[1:len(literal)-1], length, isUnicode)
	}
	return
}

func (self *Parser) CheckNewLine(quote rune) (literal string, parsed unistring.String, err error) {
	self.ScanNewline()
	errStr := "String not terminated"
	if quote == '/' {
		errStr = "Invalid regular expression: missing /"
		self.Error(self.Index, errStr)
	}
	return "", "", errors.New(errStr)
}

func (self *Parser) ScanNewline() {
	if self.Char == '\r' {
		self.Read()
		if self.Char != '\n' {
			return
		}
	}
	self.Read()
}

func Hex2Decimal(chr byte) (value rune, ok bool) {
	{
		chr := rune(chr)
		switch {
		case '0' <= chr && chr <= '9':
			return chr - '0', true
		case 'a' <= chr && chr <= 'f':
			return chr - 'a' + 10, true
		case 'A' <= chr && chr <= 'F':
			return chr - 'A' + 10, true
		}
		return
	}
}

func ParseNumberLiteral(literal string) (value interface{}, err error) {
	value, err = strconv.ParseInt(literal, 0, 64)
	if err == nil {
		return
	}

	parseIntErr := err
	value, err = strconv.ParseFloat(literal, 64)
	if err == nil {
		return
	} else if err.(*strconv.NumError).Err == strconv.ErrRange {
		return value, nil
	}

	err = parseIntErr

	if err.(*strconv.NumError).Err == strconv.ErrRange {
		if len(literal) > 2 && literal[0] == '0' && (literal[1] == 'X' || literal[1] == 'x') {
			var value float64
			literal = literal[2:]
			for _, chr := range literal {
				digit := DigitValue(chr)
				if digit >= 16 {
					return nil, errors.New("Illegal numeric literal")
				}
				value = value*16 + float64(digit)
			}
			return value, nil
		}
	}
	return nil, errors.New("Illegal numeric literal")
}

func ParseStringLiteral1(literal string, length int, unicode bool) (unistring.String, error) {
	var sb strings.Builder
	var chars []uint16
	if unicode {
		chars = make([]uint16, 1, length+1)
		chars[0] = unistring.BOM
	} else {
		sb.Grow(length)
	}
	str := literal
	for len(str) > 0 {
		switch chr := str[0]; {
		case chr >= utf8.RuneSelf:
			chr, size := utf8.DecodeRuneInString(str)
			if chr <= 0xFFFF {
				chars = append(chars, uint16(chr))
			} else {
				first, second := utf16.EncodeRune(chr)
				chars = append(chars, uint16(first), uint16(second))
			}
			str = str[size:]
			continue
		case chr != '\\':
			if unicode {
				chars = append(chars, uint16(chr))
			} else {
				sb.WriteByte(chr)
			}
			str = str[1:]
			continue
		}

		if len(str) <= 1 {
			panic("len(str) <= 1")
		}
		chr := str[1]
		var value rune
		if chr >= utf8.RuneSelf {
			str = str[1:]
			var size int
			value, size = utf8.DecodeRuneInString(str)
			str = str[size:]
			if value == '\u2028' || value == '\u2029' {
				continue
			}
		} else {
			str = str[2:]
			switch chr {
			case 'b':
				value = '\b'
			case 'f':
				value = '\f'
			case 'n':
				value = '\n'
			case 'r':
				value = '\r'
			case 't':
				value = '\t'
			case 'v':
				value = '\v'
			case 'x', 'u':
				size := 0
				switch chr {
				case 'x':
					size = 2
				case 'u':
					size = 4
				}
				if len(str) < size {
					return "", fmt.Errorf("Invalid escape: \\%s: len(%q) != %d", string(chr), str, size)
				}
				for j := 0; j < size; j++ {
					decimal, ok := Hex2Decimal(str[j])
					if !ok {
						return "", fmt.Errorf("Invalid escape: \\%s: %q", string(chr), str[:size])
					}
					value = value<<4 | decimal
				}
				str = str[size:]
				if chr == 'x' {
					break
				}
				if value > utf8.MaxRune {
					panic("value > utf8.MaxRune")
				}
			case '0':
				if len(str) == 0 || '0' > str[0] || str[0] > '7' {
					value = 0
					break
				}
				fallthrough
			case '1', '2', '3', '4', '5', '6', '7':
				value = rune(chr) - '0'
				j := 0
				for ; j < 2; j++ {
					if len(str) < j+1 {
						break
					}
					chr := str[j]
					if '0' > chr || chr > '7' {
						break
					}
					decimal := rune(str[j]) - '0'
					value = (value << 3) | decimal
				}
				str = str[j:]
			case '\\':
				value = '\\'
			case '\'', '"':
				value = rune(chr)
			case '\r':
				if len(str) > 0 {
					if str[0] == '\n' {
						str = str[1:]
					}
				}
				fallthrough
			case '\n':
				continue
			default:
				value = rune(chr)
			}
		}
		if unicode {
			if value <= 0xFFFF {
				chars = append(chars, uint16(value))
			} else {
				first, second := utf16.EncodeRune(value)
				chars = append(chars, uint16(first), uint16(second))
			}
		} else {
			if value >= utf8.RuneSelf {
				return "", fmt.Errorf("Unexpected unicode character")
			}
			sb.WriteByte(byte(value))
		}
	}
	if unicode {
		if len(chars) != length+1 {
			panic(fmt.Errorf("Unexpected unicode length while parsing '%s'", literal))
		}
		return unistring.FromUtf16(chars), nil
	}
	if sb.Len() != length {
		panic(fmt.Errorf("Unexpected length while parsing '%s'", literal))
	}
	return unistring.String(sb.String()), nil
}

func (self *Parser) ScanNumericLiteral(decimalPoint bool) (token.Token, string) {
	offset := self.CharOffset
	tkn := token.NUMBER
	if decimalPoint {
		offset--
		self.ScanMantissa(10)
		goto exponent
	}

	if self.Char == '0' {
		offset := self.CharOffset
		self.Read()
		if self.Char == 'x' || self.Char == 'X' {
			self.Read()
			if IsDigit(self.Char, 16) {
				self.Read()
			} else {
				return token.ILLEGAL, self.Template[offset:self.CharOffset]
			}
			self.ScanMantissa(16)

			if self.CharOffset-offset <= 2 {
				self.Error(self.Index, "Illegal hexadecimal number")
			}

			goto hexadecimal
		} else if self.Char == '.' {
			goto float
		} else {
			if self.Char == 'e' || self.Char == 'E' {
				goto exponent
			}
			self.ScanMantissa(8)
			if self.Char == '8' || self.Char == '9' {
				return token.ILLEGAL, self.Template[offset:self.CharOffset]
			}
			goto octal
		}
	}
	self.ScanMantissa(10)

float:
	if self.Char == '.' {
		self.Read()
		self.ScanMantissa(10)
	}
exponent:
	if self.Char == 'e' || self.Char == 'E' {
		self.Read()
		if self.Char == '-' || self.Char == '+' {
			self.Read()
		}
		if IsDecimalDigit(self.Char) {
			self.Read()
			self.ScanMantissa(10)
		} else {
			return token.ILLEGAL, self.Template[offset:self.CharOffset]
		}
	}
hexadecimal:
octal:
	if IsIdentifierStart(self.Char) || IsDecimalDigit(self.Char) {
		return token.ILLEGAL, self.Template[offset:self.CharOffset]
	}

	return tkn, self.Template[offset:self.CharOffset]
}
