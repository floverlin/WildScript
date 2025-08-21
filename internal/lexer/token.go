package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENTIFIER TokenType = "IDENTIFIER"
	LET        TokenType = "LET"

	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	FUNCTION TokenType = "FUNCTION"
	LAMBDA   TokenType = "LAMBDA"
	METHOD   TokenType = "METHOD"

	IF   TokenType = "IF"
	ELIF TokenType = "ELIF"
	ELSE TokenType = "ELSE"
	THEN TokenType = "THEN"

	FOR    TokenType = "FOR"
	IN     TokenType = "IN"
	WHILE  TokenType = "WHILE"
	DO     TokenType = "DO"
	REPEAT TokenType = "REPEAT"
	UNTIL  TokenType = "UNTIL"

	RETURN   TokenType = "RETURN"
	CONTINUE TokenType = "CONTINUE"
	BREAK    TokenType = "BREAK"

	IMPORT TokenType = "IMPORT"
	EXPORT TokenType = "EXPORT"

	AND TokenType = "AND"
	OR  TokenType = "OR"
	NOT TokenType = "NOT"

	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"

	NIL TokenType = "NIL"

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

	LARROW TokenType = "<-"
	RARROW TokenType = "->"
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
}

var dual = map[string]TokenType{
	"//": INT_DIVIDE,

	"==": EQUAL,
	"!=": NOT_EQUAL,
	"<=": LESS_EQ,
	">=": GREATER_EQ,
}

var specialIdents = map[string]TokenType{
	"let": LET,

	"function": FUNCTION,
	"lambda":   LAMBDA,
	"method":   METHOD,

	"if":   IF,
	"elif": ELIF,
	"else": ELSE,
	"then": THEN,

	"for":    FOR,
	"in":     IN,
	"while":  WHILE,
	"do":     DO,
	"repeat": REPEAT,
	"until":  UNTIL,

	"return":   RETURN,
	"continue": CONTINUE,
	"break":    BREAK,

	"import": IMPORT,
	"export": EXPORT,

	"and": AND,
	"or":  OR,
	"not": NOT,

	"true":  TRUE,
	"false": FALSE,

	"nil": NIL,
}
