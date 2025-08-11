package parser

import (
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
		stmt := p.parseStatement() // include ;
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken() // to next statement or EOF
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.IDENT:
		if p.peekToken.Type == lexer.ASSIGN {
			return p.parseVarStatement()
		} else {
			return p.parseExpressionStatement()
		}
	case lexer.SEMICOLON:
		return nil
	default:
		return p.parseExpressionStatement()
	}
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
