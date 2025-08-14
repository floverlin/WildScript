package parser

import "wildscript/internal/lexer"

const (
	LOWEST = iota

	TERNARY_LOOP // ? {loop}

	LOGICAL_OR  // ||
	LOGICAL_AND // &&

	COMPARISON // == != < > <= >=

	SUM     // + -
	PRODUCT // * / // %

	PREFIX // !

	POW // ^

	CALL // (call) [index] .
)

var precedences = map[lexer.TokenType]int{
	lexer.LBRACE:   TERNARY_LOOP,
	lexer.QUESTION: TERNARY_LOOP,

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

	lexer.POW: POW,

	lexer.NOT: PREFIX,
	lexer.AND: LOGICAL_AND,
	lexer.OR:  LOGICAL_OR,

	lexer.LPAREN:   CALL,
	lexer.DOT:      CALL,
	lexer.LBRACKET: CALL,
}
