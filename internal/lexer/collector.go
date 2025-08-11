package lexer

type Collector struct {
	lexer    *Lexer
	tokens   []Token
	illegals []Token
	pos      int
}

func NewCollector(lexer *Lexer) *Collector {
	c := &Collector{
		lexer:  lexer,
		tokens: []Token{},
		pos: 0,
	}
	c.collect()
	return c
}

func (c *Collector) NextToken() Token {
	if c.pos >= len(c.tokens) {
		panic("tokenizer: no more tokens")
	}
	token := c.tokens[c.pos]
	c.pos++
	return token
}

func (c *Collector) Tokens() []Token {
	return c.tokens[:]
}

func (c *Collector) Illegals() []Token {
	return c.illegals[:]
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
