package parser

import (
	"fmt"
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
	"wildscript/internal/lib"
)

type Tokenizer interface {
	NextToken() lexer.Token
}

type Parser struct {
	lexer     Tokenizer
	curToken  lexer.Token
	peekToken lexer.Token
}

func New(lexer Tokenizer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken() // fill peek
	p.nextToken() // fill cur
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for {
		stmt := p.parseStatement() // include ; or EOF
		program.Statements = append(program.Statements, stmt)

		if p.curToken.Type == lexer.EOF {
			break
		}

		p.nextToken() // to statement
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == lexer.SEMICOLON ||
		p.curToken.Type == lexer.EOF ||
		p.curToken.Type == lexer.RBRACE {
		return newNilStatement(p.curToken)
	}

	var stmt ast.Statement

	switch p.curToken.Type {
	case lexer.LET:
		letStmt := &ast.LetStatement{
			Token: p.curToken,
		}
		if p.peekToken.Type != lexer.IDENTIFIER {
			p.expected("identifier")
		}
		p.nextToken() // to ident
		letStmt.Left = p.parseIdentifier()
		if p.peekToken.Type != lexer.ASSIGN {
			p.expected("=")
		}
		p.nextToken() // to =
		p.nextToken() // to expr
		letStmt.Right = p.parseExpression(LOWEST)
		stmt = letStmt
	case lexer.IMPORT:
		importStmt := &ast.ImportStatement{Token: p.curToken}
		if p.peekToken.Type != lexer.IDENTIFIER {
			p.expected("expected module identifier")
		}
		p.nextToken() // to ident
		importStmt.Module = p.parseIdentifier()
		stmt = importStmt
	case lexer.EXPORT:
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

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func die(token lexer.Token, text string, args ...any) {
	text = fmt.Sprintf("[parser] %s", text)
	lib.Die(token, text, args)
}

func (p *Parser) expected(text string) {
	text = "expected " + text
	die(p.peekToken, text)
}

func newNilBlockExpression(token lexer.Token) *ast.BlockExpression {
	return &ast.BlockExpression{
		Token: token,
		Statements: []ast.Statement{
			newNilStatement(token),
		},
	}
}

func newNilStatement(token lexer.Token) ast.Statement {
	return &ast.ExpressionStatement{
		Token: token,
		Expression: &ast.NilLiteral{
			Token: token,
		},
	}
}
