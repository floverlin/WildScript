package parser

import (
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == lexer.SEMICOLON ||
		p.curToken.Type == lexer.EOF ||
		p.curToken.Type == lexer.RBRACE {
		return newNilStatement(p.curToken)
	}

	var stmt ast.Statement

	switch p.curToken.Type {
	case lexer.WHILE:
		stmt = p.parseWhileStatement()
	case lexer.REPEAT:
		stmt = p.parseRepeatStatement()
	case lexer.LET:
		letStmt := &ast.LetStatement{
			Token: p.curToken,
		}
		if p.peekToken.Type != lexer.IDENTIFIER {
			p.expected("identifier")
		}
		p.nextToken() // to ident
		letStmt.Left = p.parseIdentifier()

		if p.peekToken.Type == lexer.SEMICOLON ||
			p.peekToken.Type == lexer.EOF ||
			p.peekToken.Type == lexer.RBRACE {
			letStmt.Right = &ast.NilLiteral{Token: p.curToken}
			stmt = letStmt
		} else {
			if p.peekToken.Type != lexer.ASSIGN {
				p.expected("=")
			}
			p.nextToken() // to =
			p.nextToken() // to expr
			letStmt.Right = p.parseExpression(LOWEST)
			stmt = letStmt
		}
	case lexer.IMPORT:
		importStmt := &ast.ImportStatement{Token: p.curToken}
		if p.peekToken.Type != lexer.IDENTIFIER {
			p.expected("expected module identifier")
		}
		p.nextToken() // to ident
		importStmt.Module = append(importStmt.Module, p.parseIdentifier())
		for p.peekToken.Type == lexer.DOT {
			p.nextToken() //to .
			p.nextToken() // to ident
			importStmt.Module = append(importStmt.Module, p.parseIdentifier())
		}
		stmt = importStmt
	case lexer.EXPORT:
		exportStmt := &ast.ExportStatement{Token: p.curToken}
		p.nextToken() // to exptr
		exportStmt.Value = p.parseExpression(LOWEST)
		stmt = exportStmt
	case lexer.RETURN:
		if p.peekToken.Type == lexer.SEMICOLON ||
			p.peekToken.Type == lexer.RBRACE ||
			p.peekToken.Type == lexer.EOF {
			stmt = &ast.ReturnStatement{
				Token: p.curToken,
				Value: &ast.NilLiteral{
					Token: p.peekToken,
				},
			}
		} else {
			returnStmt := &ast.ReturnStatement{
				Token: p.curToken,
			}
			p.nextToken() // to expr
			returnStmt.Value = p.parseExpression(LOWEST)
			stmt = returnStmt
		}
	case lexer.BREAK:
		stmt = &ast.BreakStatement{
			Token: p.curToken,
		}
	case lexer.CONTINUE:
		stmt = &ast.ContinueStatement{
			Token: p.curToken,
		}
	case lexer.FUNCTION:
		if p.peekToken.Type != lexer.IDENTIFIER {
			p.expected("function identifier")
		}
		funcStmt := &ast.FunctionStatement{
			Token: p.curToken,
		}
		p.nextToken() // to ident
		funcStmt.Identifier = p.parseIdentifier()
		if p.peekToken.Type != lexer.LPAREN {
			p.expected("(")
		}
		p.nextToken() // to (
		funcStmt.Function = p.parseFunctionLiteral(ast.FUNCTION)
		stmt = funcStmt
	}

	if stmt == nil {
		token := p.curToken
		expr := p.parseExpression(LOWEST) // not include ; or EOF
		if p.peekToken.Type == lexer.ASSIGN {
			p.nextToken() // to =
			assignStmt := &ast.AssignStatement{
				Token: p.curToken,
				Left:  expr,
			}
			p.nextToken()                                // to right expr
			assignStmt.Right = p.parseExpression(LOWEST) // not include ; or EOF
			stmt = assignStmt
		} else {
			stmt = &ast.ExpressionStatement{
				Token:      token,
				Expression: expr,
			}
		}
	}

	if stmt == nil {
		die(p.curToken, "nil statement")
	}

	if p.peekToken.Type != lexer.SEMICOLON &&
		p.peekToken.Type != lexer.EOF &&
		p.peekToken.Type != lexer.RBRACE {
		p.expected("; or }")
	}

	p.nextToken() // to ; or EOF
	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{
		Token: p.curToken,
	}
	p.nextToken() // to cond
	stmt.If = p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.DO {
		p.expected("do")
	}

	p.nextToken() // to do
	if p.peekToken.Type != lexer.LBRACE {
		p.expected("{")
	}
	p.nextToken() // to {
	stmt.Loop = p.parseBlockExpression()

	return stmt
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	stmt := &ast.RepeatStatement{
		Token: p.curToken,
	}
	if p.peekToken.Type != lexer.LBRACE {
		p.expected("{")
	}
	p.nextToken() // to {
	stmt.Loop = p.parseBlockExpression()

	if p.peekToken.Type != lexer.UNTIL {
		p.expected("until")
	}
	p.nextToken() // to until
	p.nextToken() // to cond
	stmt.Until = p.parseExpression(LOWEST)

	return stmt
}
