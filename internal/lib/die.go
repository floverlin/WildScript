package lib

import (
	"fmt"
	"wildscript/internal/lexer"
)

func Die(token lexer.Token, text string, args ...any) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}
	text = fmt.Sprintf(
		"%s at line %d column %d",
		text,
		token.Line,
		token.Column,
	)

	panic(text)
}
