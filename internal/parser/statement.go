package parser

import (
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
)

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{
		Token: p.curToken,
		Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	p.nextToken() // to =
	p.nextToken() // to expt

	stmt.Value = p.parseExpression(LOWEST) // not include ;

	if p.peekToken.Type != lexer.SEMICOLON &&
		p.peekToken.Type != lexer.EOF {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected ; or EOF",
			),
		)
		return nil
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST), // not include ;
	}

	if p.peekToken.Type != lexer.SEMICOLON &&
		p.peekToken.Type != lexer.EOF {
		p.errors = append(
			p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected ; or EOF",
			),
		)
		return nil
	}

	return stmt
}
