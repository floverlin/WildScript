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
	lib.Die(token, text, args...)
}

func (p *Parser) expected(text string) {
	text = fmt.Sprintf("expected %s", text)
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
