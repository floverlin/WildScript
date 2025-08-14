package parser

import (
	"strconv"
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
)

const (
	NONE  = "none"
	RUNE  = "rune"
	OUTER = "outer"
)

// not include ; or EOF
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var expr ast.Expression

	switch p.curToken.Type {
	case lexer.NOT:
		expr = p.parsePrefixExpression()
	case lexer.LBRACE:
		expr = p.parseBlockExpression() // include }
	case lexer.LPAREN:
		p.nextToken() // to expr
		expr = p.parseExpression(LOWEST)
		if p.peekToken.Type != lexer.RPAREN {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.peekToken.Line,
					p.peekToken.Column,
					"expected )",
				),
			)
			return nil
		}
		p.nextToken() // to )

	case lexer.IDENT:
		expr = p.parseIdentifier(NONE)
	case lexer.AMPER:
		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.peekToken.Line,
					p.peekToken.Column,
					"expected identifier after &",
				),
			)
			return nil
		}
		p.nextToken() // to identifier
		expr = p.parseIdentifier(OUTER)
	case lexer.DOG:
		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.peekToken.Line,
					p.peekToken.Column,
					"expected identifier after @",
				),
			)
			return nil
		}
		p.nextToken() // to identifier
		expr = p.parseIdentifier(RUNE)

	case lexer.FN:
		expr = p.parseFuncLiteral()
	case lexer.LBRACKET:
		expr = p.parseListLiteral()
	case lexer.NUMBER:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.curToken.Line,
					p.curToken.Column,
					"could not parse %s as float",
					p.curToken.Literal,
				),
			)
			return nil
		}
		expr = &ast.FloatLiteral{Token: p.curToken, Value: value}
	case lexer.STRING:
		expr = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case lexer.TRUE:
		expr = &ast.BooleanLiteral{Token: p.curToken, Value: true}
	case lexer.FALSE:
		expr = &ast.BooleanLiteral{Token: p.curToken, Value: false}
	case lexer.NIL:
		expr = &ast.NilLiteral{Token: p.curToken}
	default:
		p.errors = append(
			p.errors,
			logger.Slog(
				p.curToken.Line,
				p.curToken.Column,
				"unexpected token: %s",
				p.curToken.Literal,
			),
		)
		return nil
	}

	for p.peekToken.Type != lexer.SEMICOLON &&
		p.peekToken.Type != lexer.EOF &&
		precedence < p.peekPrecedence() {

		p.nextToken() // to op

		switch p.curToken.Type {
		case lexer.DOT:
			expr = p.parsePropertyAccessExpression(expr)
		case lexer.LBRACKET:
			expr = p.parseBracketExpression(expr)
		case lexer.LPAREN:
			expr = p.parseCallExpression(expr)
		case lexer.LBRACE:
			expr = p.parseForExpression(expr)
		case lexer.QUESTION:
			expr = p.parseConditionExpression(expr)
		default:
			expr = p.parseInfixExpression(expr)
		}
	}

	return expr
}

func (p *Parser) parseIdentifier(identType string) *ast.Identifier {
	var isOuter, isRune bool

	switch identType {
	case "rune":
		isRune = true
	case "outer":
		isOuter = true
	}

	expr := &ast.Identifier{
		Token:   p.curToken,
		Value:   p.curToken.Literal,
		IsOuter: isOuter,
		IsRune:  isRune,
	}

	return expr
}

func (p *Parser) parseConditionExpression(
	cond ast.Expression,
) *ast.ConditionExpression {
	expr := &ast.ConditionExpression{
		Token:     p.curToken,
		Condition: cond,
	}

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected {",
			),
		)
		return nil
	}

	p.nextToken() // to {
	expr.Consequence = p.parseBlockExpression()

	if p.peekToken.Type != lexer.COLON {
		expr.Alternative = newNilBlockExpression(p.peekToken)
		return expr
	}

	p.nextToken() // to :
	p.nextToken() // to expr
	expr.Alternative = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseForExpression(
	cond ast.Expression,
) *ast.ForExpression {
	expr := &ast.ForExpression{
		Token:     p.curToken,
		Condition: cond,
	}

	expr.Body = p.parseBlockExpression() // include }

	return expr
}

func (p *Parser) parsePropertyAccessExpression(
	left ast.Expression,
) *ast.PropertyAccessExpression {
	expr := &ast.PropertyAccessExpression{
		Token:  p.curToken,
		Object: left,
	}
	p.nextToken() // to prop

	var prop *ast.Identifier
	switch p.curToken.Type {
	case lexer.IDENT:
		prop = p.parseIdentifier(NONE)
	case lexer.DOG:
		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.peekToken.Line,
					p.peekToken.Column,
					"expected identifier after @",
				),
			)
			return nil
		}
		p.nextToken() // to identifier
		prop = p.parseIdentifier(RUNE)
	default:
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected identifier after .",
			),
		)
		return nil
	}

	expr.Property = prop
	return expr
}

func (p *Parser) parseBracketExpression(
	left ast.Expression,
) ast.Expression {
	var expr ast.Expression
	token := p.curToken

	p.nextToken() // to index
	firstIndex := p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // to :
		p.nextToken() // to index

		secondIndex := p.parseExpression(LOWEST)
		expr = &ast.SliceExpression{
			Token: token,
			Left:  left,
			Start: firstIndex,
			End:   secondIndex,
		}
	} else {
		expr = &ast.IndexExpression{
			Token: token,
			Left:  left,
			Index: firstIndex,
		}
	}

	if p.peekToken.Type != lexer.RBRACKET {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected ]",
			),
		)
		return nil
	}

	p.nextToken() // to ]
	return expr
}

func (p *Parser) parsePrefixExpression() *ast.PrefixExpression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken() // to expr
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) *ast.InfixExpression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken() // to right expr
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseCallExpression(
	function ast.Expression,
) *ast.CallExpression {
	expr := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	expr.Arguments = p.parseCallArguments() // include )
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken() // to )
		return args
	}

	p.nextToken() // to args
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // to ,
		p.nextToken() // to expr
		args = append(args, p.parseExpression(LOWEST))
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected )",
			),
		)
		return nil
	}

	p.nextToken() // to )
	return args
}

func (p *Parser) parseFuncLiteral() *ast.FuncLiteral {
	function := &ast.FuncLiteral{
		Token: p.curToken,
	}

	if p.peekToken.Type != lexer.LPAREN {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected (",
			),
		)
		return nil
	}

	p.nextToken()                                 // to (
	function.Parameters = p.parseFuncParameters() // include )

	if p.peekToken.Type != lexer.LBRACE {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected {",
			),
		)
		return nil
	}

	p.nextToken() // to {
	function.Body = p.parseBlockExpression()

	return function
}

func (p *Parser) parseFuncParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken() // to )
		return params
	}

	p.nextToken() // to ident
	params = append(params, p.parseIdentifier(NONE))

	for {
		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken() // to ,
		p.nextToken() // to ident
		params = append(params, p.parseIdentifier(NONE))
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected )",
			),
		)
		return nil
	}

	p.nextToken() // to )
	return params
}

// include {}
func (p *Parser) parseBlockExpression() *ast.BlockExpression {
	block := &ast.BlockExpression{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // to statement

	for {
		stmt := p.parseStatement() // include ; and }
		block.Statements = append(block.Statements, stmt)

		if p.curToken.Type == lexer.RBRACE {
			break
		}

		p.nextToken() // to statement
	}

	return block
}

func (p *Parser) parseListLiteral() *ast.ListLiteral {
	lit := &ast.ListLiteral{
		Token: p.curToken,
	}
	elems := []ast.Expression{}

	if p.peekToken.Type == lexer.RBRACKET {
		p.nextToken() // to ]
		lit.Elements = elems
		return lit
	}

	p.nextToken() // to elem
	elems = append(elems, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // to ,
		p.nextToken() // to elem

		elems = append(elems, p.parseExpression(LOWEST))
	}

	if p.peekToken.Type != lexer.RBRACKET {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected ]",
			),
		)
		return nil
	}

	p.nextToken() // to ]
	lit.Elements = elems
	return lit
}
