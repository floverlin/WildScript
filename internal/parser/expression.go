package parser

import (
	"fmt"
	"strconv"
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
)

// not include ;
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var expr ast.Expression

	switch p.curToken.Type {
	case lexer.NOT:
		expr = p.parsePrefixExpression()
	case lexer.LBRACE:
		p.nextToken()                   // to block statement
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
		expr = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if p.peekToken.Type == lexer.LPAREN {
			p.nextToken() // to (
			expr = p.parseCallExpression(expr)
		}
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
		p.nextToken() // to ident
		expr = &ast.Identifier{
			Token:   p.curToken,
			Value:   p.curToken.Literal,
			IsOuter: true,
		}
		if p.peekToken.Type == lexer.LPAREN {
			p.nextToken() // to (
			expr = p.parseCallExpression(expr)
		}
	case lexer.NUMBER:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(
				p.errors,
				logger.Slog(
					p.curToken.Line,
					p.curToken.Column,
					"could not parse: %s as float",
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

	for p.peekToken.Type != lexer.SEMICOLON {
		if precedence < p.peekPrecedence() {
			p.nextToken() // to operator
			if p.curToken.Type == lexer.QUESTION {
				expr = p.parseConditionExpression(expr)
			} else {
				expr = p.parseInfixExpression(expr)
				fmt.Println(expr)
			}
			if p.peekToken.Type == lexer.LBRACE {
				p.nextToken()
				expr = p.parseWhileExpression(expr)
			}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) parseConditionExpression(
	cond ast.Expression,
) *ast.ConditionExpression {
	expr := &ast.ConditionExpression{
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
	p.nextToken() // to block statement

	expr.Consequence = p.parseBlockExpression()

	if p.peekToken.Type != lexer.COLON {
		expr.Alternative = nilBlock(p.curToken)
		return expr
	}

	p.nextToken() // to :
	p.nextToken() // to expr
	expr.Alternative = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseWhileExpression(
	cond ast.Expression,
) *ast.WhileExpression {
	expr := &ast.WhileExpression{
		Token:     p.curToken,
		Condition: cond,
	}

	p.nextToken() // to block statement
	expr.Body = p.parseBlockExpression()

	return expr
}

func nilBlock(token lexer.Token) *ast.BlockExpression {
	return &ast.BlockExpression{
		Token: token,
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token,
				Expression: &ast.NilLiteral{
					Token: token,
				},
			},
		},
	}
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

func (p *Parser) parseCallExpression(function ast.Expression) *ast.CallExpression {
	expr := &ast.CallExpression{Token: p.curToken, Function: function}
	expr.Arguments = p.parseCallArguments()
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

func (p *Parser) parseBlockExpression() *ast.BlockExpression {
	block := &ast.BlockExpression{Token: p.curToken}
	block.Statements = []ast.Statement{}

	if p.curToken.Type == lexer.RBRACE {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		return block
	}

	for p.curToken.Type != lexer.RBRACE {
		if p.curToken.Type == lexer.SEMICOLON {
			p.nextToken() // to statement or RBRACE
		}
		stmt := p.parseStatement() // not include ;
		block.Statements = append(block.Statements, stmt)
		if p.curToken.Type == lexer.RBRACE {
			break
		}
		p.nextToken() // to ; or RBRACE
	}

	return block
}
