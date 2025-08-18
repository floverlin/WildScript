package parser

import "wildscript/internal/lexer"

const (
	LOWEST = iota
	IF_LOOP
	LOGICAL_OR
	LOGICAL_AND
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	POW
	CALL
	HIGHEST
)

var precedences = map[lexer.TokenType]int{
	lexer.IF:    IF_LOOP,
	lexer.FOR:   IF_LOOP,
	lexer.WHILE: IF_LOOP,
	lexer.UNTIL: IF_LOOP,

	lexer.OR:  LOGICAL_OR,
	lexer.AND: LOGICAL_AND,

	lexer.EQUAL:      COMPARISON,
	lexer.NOT_EQUAL:  COMPARISON,
	lexer.LESS:       COMPARISON,
	lexer.GREATER:    COMPARISON,
	lexer.LESS_EQ:    COMPARISON,
	lexer.GREATER_EQ: COMPARISON,

	lexer.PLUS:  SUM,
	lexer.MINUS: SUM,

	lexer.MULTIPLY:   PRODUCT,
	lexer.DIVIDE:     PRODUCT,
	lexer.INT_DIVIDE: PRODUCT,
	lexer.MOD:        PRODUCT,

	lexer.NOT: PREFIX,

	lexer.POW: POW,

	lexer.DOT:      CALL,
	lexer.LPAREN:   CALL,
	lexer.LBRACKET: CALL,
	lexer.LBRACE:   CALL,
}
