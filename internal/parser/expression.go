package parser

import (
	"strconv"
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
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
		expr = p.parseIdentifier(false)
	case lexer.AMPER:
		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.peekToken.Line,
					p.peekToken.Column,
					"expected identifier",
				),
			)
			return nil
		}
		p.nextToken() // to identifier
		expr = p.parseIdentifier(true)

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

	for {

		if p.peekToken.Type == lexer.SEMICOLON ||
			p.peekToken.Type == lexer.EOF {
			break
		}

		nextPrec := p.peekPrecedence()
		if p.peekToken.Type == lexer.LBRACE {
			nextPrec = TERNARY_LOOP
		} else if p.peekToken.Type == lexer.QUESTION {
			nextPrec = TERNARY_LOOP
		}

		if precedence >= nextPrec {
			break
		}

		if p.peekToken.Type == lexer.QUESTION {
			p.nextToken() // to ?
			expr = p.parseConditionExpression(expr)
			continue
		}

		if p.peekToken.Type == lexer.LBRACE {
			p.nextToken() // to {
			expr = p.parseForExpression(expr)
			continue
		}

		p.nextToken() // to operator
		expr = p.parseInfixExpression(expr)
	}

	return expr
}

// CONTAINS FUNC // TODO REMOVE FUNC CALL
func (p *Parser) parseIdentifier(isOuter bool) ast.Expression {
	var expr ast.Expression = &ast.Identifier{
		Token:   p.curToken,
		Value:   p.curToken.Literal,
		IsOuter: isOuter,
	}
	if p.peekToken.Type == lexer.LPAREN { // TODO AFTER FUNCS
		p.nextToken() // to (
		expr = p.parseCallExpression(expr)
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
