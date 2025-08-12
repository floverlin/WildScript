package parser

import (
	"slices"
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
)

type Parser struct {
	lexer     lexer.Tokenizer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(lexer lexer.Tokenizer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	if p.curToken.Type != lexer.EOF {
		p.peekToken = p.lexer.NextToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.SEMICOLON {
			p.nextToken() // to statement or EOF
		}
		stmt := p.parseStatement() // not include ;
		program.Statements = append(program.Statements, stmt)
		p.nextToken() // to ; or EOF
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if slices.Contains(
		[]lexer.TokenType{lexer.SEMICOLON, lexer.EOF, lexer.RBRACE},
		p.curToken.Type,
	) {
		return &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: &ast.NilLiteral{Token: p.curToken},
		}
	}

	var stmt ast.Statement
	expr := p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // to ==
		p.nextToken() // to right expr
		assign := &ast.AssignStatement{
			Token: p.curToken,
			Left:  expr,
		}
		right := p.parseExpression(LOWEST)
		assign.Right = right
		stmt = assign
	} else {
		stmt = &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: expr,
		}
	}

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

func (p *Parser) Errors() []string {
	return p.errors
}
