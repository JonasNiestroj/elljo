package token

import "strconv"

type Token int

const (
	_ Token = iota

	ILLEGAL
	EOF
	COMMENT
	KEYWORD

	STRING
	BOOLEAN
	NULL
	NUMBER
	IDENTIFIER

	PLUS
	MINUS
	MULTIPLY
	SLASH
	REMAINDER

	AND
	OR
	EXCLUSIVE_OR
	SHIFT_LEFT
	SHIFT_RIGHT
	UNSIGNED_SHIFT_RIGHT

	ADD_ASSIGN
	SUBTRACT_ASSIGN
	MULTIPLY_ASSIGN
	QUOTIENT_ASSIGN
	REMAINDER_ASSIGN

	AND_ASSIGN
	OR_ASSIGN
	EXCLUSIVE_OR_ASSIGN
	SHIFT_LEFT_ASSIGN
	SHIFT_RIGHT_ASSIGN
	UNSIGNED_SHIFT_RIGHT_ASSIGN

	LOGICAL_AND
	LOGICAL_OR
	INCREMENT
	DECREMENT

	EQUAL
	STRICT_EQUAL
	LESS
	GREATER
	ASSIGN
	NOT

	BITWISE_NOT

	NOT_EQUAL
	STRICT_NOT_EQUAL
	LESS_OR_EQUAL
	GREATER_OR_EQUAL

	LEFT_PARENTHESIS
	LEFT_BRACKET
	LEFT_BRACE
	COMMA
	PERIOD
	Ellipsis

	RIGHT_PARENTHESIS
	RIGHT_BRACKET
	RIGHT_BRACE
	SEMICOLON
	COLON
	QUESTION_MARK

	IF
	IN
	OF
	DO

	VAR
	LET
	CONST
	FOR
	NEW
	TRY

	THIS
	ELSE
	CASE
	VOID
	WITH

	WHILE
	BREAK
	CATCH
	THROW

	RETURN
	TYPEOF
	DELETE
	SWITCH

	DEFAULT
	FINALLY

	FUNCTION
	CONTINUE
	DEBUGGER

	INSTANCEOF
	IMPORT
	IMPORTFROM
	lastKeyword
)

var tokenMap = [...]string{
	ILLEGAL:                     "ILLEGAL",
	EOF:                         "EOF",
	COMMENT:                     "COMMENT",
	KEYWORD:                     "KEYWORD",
	STRING:                      "STRING",
	BOOLEAN:                     "BOOLEAN",
	NULL:                        "NULL",
	NUMBER:                      "NUMBER",
	IDENTIFIER:                  "IDENTIFIER",
	PLUS:                        "+",
	MINUS:                       "-",
	MULTIPLY:                    "*",
	SLASH:                       "/",
	REMAINDER:                   "%",
	AND:                         "&",
	OR:                          "|",
	EXCLUSIVE_OR:                "^",
	SHIFT_LEFT:                  "<<",
	SHIFT_RIGHT:                 ">>",
	UNSIGNED_SHIFT_RIGHT:        ">>>",
	ADD_ASSIGN:                  "+=",
	SUBTRACT_ASSIGN:             "-=",
	MULTIPLY_ASSIGN:             "*=",
	QUOTIENT_ASSIGN:             "/=",
	REMAINDER_ASSIGN:            "%=",
	AND_ASSIGN:                  "&=",
	OR_ASSIGN:                   "|=",
	EXCLUSIVE_OR_ASSIGN:         "^=",
	SHIFT_LEFT_ASSIGN:           "<<=",
	SHIFT_RIGHT_ASSIGN:          ">>=",
	UNSIGNED_SHIFT_RIGHT_ASSIGN: ">>>=",
	LOGICAL_AND:                 "&&",
	LOGICAL_OR:                  "||",
	INCREMENT:                   "++",
	DECREMENT:                   "--",
	EQUAL:                       "==",
	STRICT_EQUAL:                "===",
	LESS:                        "<",
	GREATER:                     ">",
	ASSIGN:                      "=",
	NOT:                         "!",
	BITWISE_NOT:                 "~",
	NOT_EQUAL:                   "!=",
	STRICT_NOT_EQUAL:            "!==",
	LESS_OR_EQUAL:               "<=",
	GREATER_OR_EQUAL:            ">=",
	LEFT_PARENTHESIS:            "(",
	LEFT_BRACKET:                "[",
	LEFT_BRACE:                  "{",
	COMMA:                       ",",
	PERIOD:                      ".",
	Ellipsis:                    "...",
	RIGHT_PARENTHESIS:           ")",
	RIGHT_BRACKET:               "]",
	RIGHT_BRACE:                 "}",
	SEMICOLON:                   ";",
	COLON:                       ":",
	QUESTION_MARK:               "?",
	IF:                          "if",
	IN:                          "in",
	OF:                          "of",
	DO:                          "do",
	VAR:                         "var",
	LET:                         "let",
	CONST:                       "const",
	FOR:                         "for",
	NEW:                         "new",
	TRY:                         "try",
	THIS:                        "this",
	ELSE:                        "else",
	CASE:                        "case",
	VOID:                        "void",
	WITH:                        "with",
	WHILE:                       "while",
	BREAK:                       "break",
	CATCH:                       "catch",
	THROW:                       "throw",
	RETURN:                      "return",
	TYPEOF:                      "typeof",
	DELETE:                      "delete",
	SWITCH:                      "switch",
	DEFAULT:                     "default",
	FINALLY:                     "finally",
	FUNCTION:                    "function",
	CONTINUE:                    "continue",
	DEBUGGER:                    "debugger",
	INSTANCEOF:                  "instanceof",
	IMPORT:                      "import",
	IMPORTFROM:                  "from",
}

type Keyword struct {
	Token Token

}

var keywordMap = map[string]Keyword{
	"if": {
		Token: IF,
	},
	"in": {
		Token: IN,
	},
	"do": {
		Token: DO,
	},
	"var": {
		Token: VAR,
	},
	"let": {
		Token: LET,
	},
	"const": {
		Token: CONST,
	},
	"for": {
		Token: FOR,
	},
	"new": {
		Token: NEW,
	},
	"try": {
		Token: TRY,
	},
	"this": {
		Token: THIS,
	},
	"else": {
		Token: ELSE,
	},
	"case": {
		Token: CASE,
	},
	"void": {
		Token: VOID,
	},
	"with": {
		Token: WITH,
	},
	"while": {
		Token: WHILE,
	},
	"break": {
		Token: BREAK,
	},
	"catch": {
		Token: CATCH,
	},
	"throw": {
		Token: THROW,
	},
	"return": {
		Token: RETURN,
	},
	"typeof": {
		Token: TYPEOF,
	},
	"delete": {
		Token: DELETE,
	},
	"switch": {
		Token: SWITCH,
	},
	"default": {
		Token: DEFAULT,
	},
	"finally": {
		Token: FINALLY,
	},
	"function": {
		Token: FUNCTION,
	},
	"continue": {
		Token: CONTINUE,
	},
	"debugger": {
		Token: DEBUGGER,
	},
	"instanceof": {
		Token: INSTANCEOF,
	},
	"class": {
		Token: KEYWORD,
	},
	"enum": {
		Token: KEYWORD,
	},
	"export": {
		Token: KEYWORD,
	},
	"extends": {
		Token: KEYWORD,
	},
	"import": {
		Token: IMPORT,
	},
	"from": {
		Token: IMPORTFROM,
	},
	"super": {
		Token: KEYWORD,
	},
	"implements": {
		Token: KEYWORD,
	},
	"interface": {
		Token: KEYWORD,
	},
	"package": {
		Token: KEYWORD,
	},
	"private": {
		Token: KEYWORD,
	},
	"protected": {
		Token: KEYWORD,
	},
	"public": {
		Token: KEYWORD,
	},
	"static": {
		Token: KEYWORD,
	},
}

func (token Token) ToString() string {
	if token == 0 {
		return "Not defined token"
	}
	if token < Token(len(tokenMap)) {
		return tokenMap[token]
	}
	// Return a generic token
	return "token(" + strconv.Itoa(int(token)) + ")"
}


func StringIsKeyword(string string) Token {
	if keyword, exists := keywordMap[string]; exists {
		return keyword.Token
	}
	return 0
}