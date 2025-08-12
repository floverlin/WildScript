package ast

import (
	"strings"
	"wildscript/internal/lexer"
)

type InfixExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	return joiner("(", ie.Left.String(), ie.Operator, ie.Right.String(), ")")
}

type PrefixExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return joiner("(", pe.Operator, pe.Right.String(), ")")
}

type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out strings.Builder
	out.WriteString(ce.Function.String() + "(")
	for idx, arg := range ce.Arguments {
		out.WriteString(arg.String())
		if idx != len(ce.Arguments)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteByte(')')
	return out.String()
}

type BlockExpression struct {
	Token      lexer.Token
	Statements []Statement
}

func (be *BlockExpression) expressionNode() {}
func (be *BlockExpression) String() string {
	var out strings.Builder
	out.WriteString("{\n")
	for _, stmt := range be.Statements {
		out.WriteString("    " + stmt.String())
		out.WriteByte('\n')
	}
	out.WriteString("}\n")
	return out.String()
}
