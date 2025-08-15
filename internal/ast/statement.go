package ast

import (
	"fmt"
	"arc/internal/lexer"
)

type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type AssignStatement struct {
	Token lexer.Token
	Left  Expression
	Right Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) String() string {
	return joiner(as.Left.String(), "=", as.Right.String())
}

type FuncStatement struct {
	Token      lexer.Token
	Identifier *Identifier
	Function   *FuncLiteral
}

func (fs *FuncStatement) statementNode() {}
func (fs *FuncStatement) String() string {
	return fmt.Sprintf(
		"%s => %s",
		fs.Identifier.String(),
		fs.Function.String(),
	)
}

type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return joiner("<-", rs.Value.String())
}

type ContinueStatement struct {
	Token lexer.Token
}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) String() string {
	return "->"
}
