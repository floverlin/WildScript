package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	NIL    TokenType = "NIL"
	FN     TokenType = "FN"

	AMPER     TokenType = "&"
	ASSIGN    TokenType = "="
	SEMICOLON TokenType = ";"
	COMMA     TokenType = ","
	QUESTION  TokenType = "?"
	COLON     TokenType = ":"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"

	PLUS       TokenType = "+"
	MINUS      TokenType = "-"
	MULTIPLY   TokenType = "*"
	DIVIDE     TokenType = "/"
	INT_DIVIDE TokenType = "//"
	MOD        TokenType = "%"
	POW        TokenType = "^"

	EQUAL      TokenType = "=="
	NOT_EQUAL  TokenType = "!="
	LESS       TokenType = "<"
	GREATER    TokenType = ">"
	LESS_EQ    TokenType = "<="
	GREATER_EQ TokenType = ">="

	AND TokenType = "&&"
	OR  TokenType = "||"
	NOT TokenType = "!"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func newToken(t TokenType, lit string, line, column int) Token {
	return Token{Type: t, Literal: lit, Line: line, Column: column}
}

var mono = map[byte]TokenType{
	'&': AMPER,
	'=': ASSIGN,
	';': SEMICOLON,
	',': COMMA,
	'?': QUESTION,
	':': COLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,

	'+': PLUS,
	'-': MINUS,
	'*': MULTIPLY,
	'/': DIVIDE,
	'%': MOD,
	'^': POW,

	'<': LESS,
	'>': GREATER,

	'!': NOT,
}

var dual = map[string]TokenType{
	"//": INT_DIVIDE,

	"==": EQUAL,
	"!=": NOT_EQUAL,
	"<=": LESS_EQ,
	">=": GREATER_EQ,

	"&&": AND,
	"||": OR,
}

var specialIdents = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"nil":   NIL,
	"fn":    FN,
}
