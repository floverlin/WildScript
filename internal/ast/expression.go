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
	out.WriteString("{ ")
	for idx, stmt := range be.Statements {
		if idx != 0 {
			out.WriteString("    ")
		}
		out.WriteString(stmt.String())
		if idx != len(be.Statements)-1 {
			out.WriteByte('\n')
		}
	}
	out.WriteString(" }")
	return out.String()
}

type ConditionExpression struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockExpression
	Alternative Expression
}

func (ce *ConditionExpression) expressionNode() {}
func (ce *ConditionExpression) String() string {
	return joiner(
		ce.Condition.String(), "?",
		ce.Consequence.String(), ":",
		ce.Alternative.String(),
	)
}

type ForExpression struct {
	Token     lexer.Token
	Condition Expression
	Body      *BlockExpression
}

func (fe *ForExpression) expressionNode() {}
func (fe *ForExpression) String() string {
	return joiner(
		fe.Condition.String(),
		fe.Body.String(),
	)
}
