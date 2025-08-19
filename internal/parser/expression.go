package parser

import (
	"strconv"
	"wildscript/internal/ast"
	"wildscript/internal/lexer"
)

// not include ; or EOF
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var expr ast.Expression

	switch p.curToken.Type {
	case lexer.NOT:
		expr = p.parsePrefixExpression()
	case lexer.MINUS:
		expr = p.parsePrefixExpression()
	case lexer.LPAREN:
		p.nextToken() // to expr
		expr = p.parseExpression(LOWEST)
		if p.peekToken.Type != lexer.RPAREN {
			die(p.peekToken, "expected )")
		}
		p.nextToken() // to )

	case lexer.IDENTIFIER:
		expr = p.parseIdentifier()

	case lexer.IF:
		expr = p.parseIfExpression()
	case lexer.ELIF:
		expr = p.parseIfExpression()

	case lexer.LAMBDA:
		p.nextToken() // to (
		expr = p.parseFunctionLiteral(ast.LAMBDA)
	case lexer.METHOD:
		p.nextToken() // to (
		expr = p.parseFunctionLiteral(ast.METHOD)
	case lexer.LBRACE:
		expr = p.parseDocumentLiteral()
	case lexer.NUMBER:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			die(
				p.curToken,
				"could not parse %s as number",
				p.curToken.Literal,
			)
		}
		expr = &ast.NumberLiteral{Token: p.curToken, Value: value}
	case lexer.STRING:
		expr = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case lexer.TRUE:
		expr = &ast.BooleanLiteral{Token: p.curToken, Value: true}
	case lexer.FALSE:
		expr = &ast.BooleanLiteral{Token: p.curToken, Value: false}
	case lexer.NIL:
		expr = &ast.NilLiteral{Token: p.curToken}
	default:
		die(
			p.curToken,
			"unexpected token: %s",
			p.curToken.Literal,
		)
	}

	for precedence < p.peekPrecedence() {

		p.nextToken() // to op

		switch p.curToken.Type {
		case lexer.DOT:
			expr = p.parsePropertyExpression(expr)
		case lexer.LBRACKET:
			expr = p.parseBracketExpression(expr)
		case lexer.LPAREN:
			expr = p.parseCallExpression(expr)
		case lexer.LBRACE:
			expr = p.parseKeyExpression(expr)
		default:
			expr = p.parseInfixExpression(expr)
		}
	}

	return expr
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	expr := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return expr
}

func (p *Parser) parseIfExpression() *ast.IfExpression {
	expr := &ast.IfExpression{
		Token: p.curToken,
	}

	p.nextToken() // to cond
	expr.If = p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.THEN {
		p.expected("then after condition")
	}
	p.nextToken() // to then

	if p.peekToken.Type != lexer.LBRACE {
		p.expected("{")
	}

	p.nextToken() // to {
	expr.Then = p.parseBlockExpression()

	if p.peekToken.Type != lexer.ELSE &&
		p.peekToken.Type != lexer.ELIF {
		expr.Else = newNilBlockExpression(p.peekToken)
		return expr
	}

	p.nextToken() // to else or elif
	if p.curToken.Type == lexer.ELIF {
		expr.Else = p.parseIfExpression()
	} else {
		p.nextToken() // to block
		expr.Else = p.parseBlockExpression()
	}

	return expr
}

func (p *Parser) parsePropertyExpression(
	left ast.Expression,
) *ast.PropertyExpression {
	expr := &ast.PropertyExpression{
		Token: p.curToken,
		Left:  left,
	}

	if p.peekToken.Type != lexer.IDENTIFIER {
		p.expected("property")
	}

	p.nextToken() // to prop
	expr.Property = p.parseIdentifier()

	return expr
}

func (p *Parser) parseKeyExpression(
	left ast.Expression,
) *ast.KeyExpression {
	expr := &ast.KeyExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.nextToken() // to key
	expr.Key = p.parseExpression(LOWEST)

	if p.peekToken.Type != lexer.RBRACE {
		p.expected("}")
	}

	p.nextToken() // to }
	return expr
}

func (p *Parser) parseBracketExpression(
	left ast.Expression,
) ast.Expression {
	var expr ast.Expression
	token := p.curToken

	p.nextToken() // to index
	firstIndex := p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // to :
		p.nextToken() // to index

		secondIndex := p.parseExpression(LOWEST)
		expr = &ast.SliceExpression{
			Token: token,
			Left:  left,
			Start: firstIndex,
			End:   secondIndex,
		}
	} else {
		expr = &ast.IndexExpression{
			Token: token,
			Left:  left,
			Index: firstIndex,
		}
	}

	if p.peekToken.Type != lexer.RBRACKET {
		p.expected("]")
	}

	p.nextToken() // to ]
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
		p.expected(")")
	}

	p.nextToken() // to )
	return args
}

// from (
func (p *Parser) parseFunctionLiteral(funcImpl ast.FunctionImplementation) *ast.FunctionLiteral {
	function := &ast.FunctionLiteral{
		Token: p.curToken,
		Impl:  funcImpl,
	}

	function.Parameters = p.parseFunctionParameters() // include )

	if p.peekToken.Type != lexer.LBRACE {
		p.expected("{")
	}

	p.nextToken() // to {
	function.Body = p.parseBlockExpression()

	return function
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken() // to )
		return params
	}

	if p.peekToken.Type != lexer.IDENTIFIER {
		p.expected("parameter")
	}

	p.nextToken() // to ident
	params = append(params, p.parseIdentifier())

	for {
		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken() // to ,
		if p.peekToken.Type != lexer.IDENTIFIER {
			p.expected("parameter")
		}
		p.nextToken() // to ident
		params = append(params, p.parseIdentifier())
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.expected(")")
	}

	p.nextToken() // to )
	return params
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

func (p *Parser) parseDocumentLiteral() *ast.DocumentLiteral {
	lit := &ast.DocumentLiteral{
		Token: p.curToken,
	}
	elems := []*ast.DocumentElement{}

	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken() // to }
		lit.Elements = elems
		return lit
	}

	p.nextToken() // to elem

	elems = append(elems, p.parseDocumentElement())

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // to ,
		p.nextToken() // to elem
		elems = append(elems, p.parseDocumentElement())
	}

	if p.peekToken.Type != lexer.RBRACE {
		p.expected("}")
	}

	p.nextToken() // to }
	lit.Elements = elems
	return lit
}

func (p *Parser) parseDocumentElement() *ast.DocumentElement {
	elem := &ast.DocumentElement{
		Token: p.curToken,
	}

	left := p.parseExpression(LOWEST)

	switch p.peekToken.Type {
	case lexer.COMMA:
		elem.Type = ast.LIST
		elem.Value = left
	case lexer.ASSIGN:
		p.nextToken() // to =
		p.nextToken() // to value
		elem.Type = ast.PROP
		elem.Key = left
		elem.Value = p.parseExpression(LOWEST)
	case lexer.COLON:
		p.nextToken() // to :
		p.nextToken() // to value
		elem.Type = ast.DICT
		elem.Key = left
		elem.Value = p.parseExpression(LOWEST)
	}

	return elem
}
