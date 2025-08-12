package ast

import "wildscript/internal/lexer"

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
