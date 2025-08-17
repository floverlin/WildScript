package parser

import (
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
)

type Parser struct {
	lexer     lexer.Tokenizer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(lexer lexer.Tokenizer) *Parser {
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
	program.Statements = []ast.Statement{}

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
	token := p.curToken

	if p.curToken.Type == lexer.RETURN {
		if p.peekToken.Type == lexer.SEMICOLON ||
			p.peekToken.Type == lexer.EOF {
			stmt = &ast.ReturnStatement{
				Token: token,
				Value: &ast.NilLiteral{
					Token: token,
				},
			}
		} else {
			p.nextToken() // to expr
			stmt = &ast.ReturnStatement{
				Token: token,
				Value: p.parseExpression(LOWEST),
			}
		}
	} else if p.curToken.Type == lexer.CONTINUE {
		stmt = &ast.ContinueStatement{
			Token: token,
		}

	} else if p.curToken.Type == lexer.USE {
		useStmt := &ast.UseStatement{Token: p.curToken}
		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(p.errors,
				logger.Slog(
					p.peekToken.Line,
					p.peekToken.Column,
					"expected module identifier",
				),
			)
			return nil
		}
		p.nextToken() // to ident
		useStmt.Name = p.parseIdentifier(NONE)
		stmt = useStmt

	} else if p.curToken.Type == lexer.FN &&
		p.peekToken.Type == lexer.IDENT {
		p.nextToken() // to ident
		funcStmt := &ast.FuncStatement{
			Token: p.curToken,
		}
		funcStmt.Identifier = p.parseIdentifier(NONE)
		p.nextToken() // to (
		funcStmt.Function = p.parseFuncLiteral()
		stmt = funcStmt
	} else {
		expr := p.parseExpression(LOWEST) // not include ; or EOF

		if p.peekToken.Type == lexer.ASSIGN {
			p.nextToken() // to =
			assign := &ast.AssignStatement{
				Token: token,
				Left:  expr,
			}

			p.nextToken()                      // to right expr
			right := p.parseExpression(LOWEST) // not include ; or EOF
			assign.Right = right

			stmt = assign
		} else {
			stmt = &ast.ExpressionStatement{
				Token:      token,
				Expression: expr,
			}
		}

	}

	if p.peekToken.Type != lexer.SEMICOLON &&
		p.peekToken.Type != lexer.EOF &&
		p.peekToken.Type != lexer.RBRACE {
		p.errors = append(p.errors,
			logger.Slog(
				p.peekToken.Line,
				p.peekToken.Column,
				"expected ;",
			),
		)
		return nil
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

func (p *Parser) Errors() []string {
	return p.errors
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
