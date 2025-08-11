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

type VarStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) String() string {
	return joiner(vs.Name.String(), "=", vs.Value.String())
}
