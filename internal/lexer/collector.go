package lexer

import "slices"

type Collector struct {
	lexer    *Lexer
	tokens   []Token
	illegals []Token
	pos      int
}

func NewCollector(lexer *Lexer) *Collector {
	c := &Collector{lexer: lexer}
	c.collect()
	return c
}

func (c *Collector) NextToken() Token {
	if c.pos >= len(c.tokens) {
		return c.tokens[len(c.tokens)-1]
	}
	token := c.tokens[c.pos]
	c.pos++
	return token
}

func (c *Collector) Tokens() []Token {
	return slices.Clone(c.tokens)
}

func (c *Collector) Illegals() []Token {
	return slices.Clone(c.illegals)
}

func (c *Collector) Reset() {
	c.pos = 0
}

func (c *Collector) collect() {
	for {
		token := c.lexer.NextToken()
		c.tokens = append(c.tokens, token)

		if token.Type == ILLEGAL {
			c.illegals = append(c.illegals, token)
		}

		if token.Type == EOF {
			break
		}
	}
}
