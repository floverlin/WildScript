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
	NEW    TokenType = "NEW"
	USE    TokenType = "USE"

	DOT       TokenType = "."
	DOG       TokenType = "@"
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
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"

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

	RETURN   TokenType = "<-"
	CONTINUE TokenType = "->"
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
	'.': DOT,
	'@': DOG,
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
	'[': LBRACKET,
	']': RBRACKET,

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

	"<-": RETURN,
	"->": CONTINUE,
}

var specialIdents = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"nil":   NIL,
	"fn":    FN,
	"new":   NEW,
	"use":   USE,
}
