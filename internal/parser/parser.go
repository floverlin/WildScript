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
	switch p.curToken.Type {
	case lexer.IDENT:
		if p.peekToken.Type == lexer.ASSIGN {
			return p.parseVarStatement() // not include ;
		} else {
			return p.parseExpressionStatement() // not include ;
		}
	case lexer.SEMICOLON: // empty statement
		return &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: &ast.NilLiteral{Token: p.curToken},
		}
	case lexer.EOF: // end of program
		return &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: &ast.NilLiteral{Token: p.curToken},
		}
	case lexer.RBRACE: // end of block
		return &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: &ast.NilLiteral{Token: p.curToken},
		}
	default:
		return p.parseExpressionStatement() // not include ;
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
