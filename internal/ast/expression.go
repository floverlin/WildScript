package ast

import (
	"fmt"
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
	return fmt.Sprintf(
		"(%s %s %s)",
		ie.Left.String(),
		ie.Operator,
		ie.Right.String(),
	)
}

type PrefixExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s %s)", pe.Operator, pe.Right.String())
}

type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var sb strings.Builder
	sb.WriteString(ce.Function.String() + "(")
	for idx, arg := range ce.Arguments {
		sb.WriteString(arg.String())
		if idx != len(ce.Arguments)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

type BlockExpression struct {
	Token      lexer.Token
	Statements []Statement
}

func (be *BlockExpression) expressionNode() {}
func (be *BlockExpression) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for idx, stmt := range be.Statements {
		sb.WriteString(stmt.String())
		if idx != len(be.Statements)-1 {
			sb.WriteString("; ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

type IfExpression struct {
	Token lexer.Token
	If    Expression
	Then  *BlockExpression
	Else  Expression
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) String() string {
	return fmt.Sprintf(
		"if %s then %s else %s",
		ie.If.String(),
		ie.Then.String(),
		ie.Else.String(),
	)
}

type ForExpression struct {
	Token lexer.Token
}

func (fe *ForExpression) expressionNode() {}
func (fe *ForExpression) String() string  { return "" }

type RepeatExpression struct {
	Token lexer.Token
}

func (re *RepeatExpression) expressionNode() {}
func (re *RepeatExpression) String() string  { return "" }

type WhileExpression struct {
	Token lexer.Token
	If    Expression
	Loop  *BlockExpression
}

func (we *WhileExpression) expressionNode() {}
func (we *WhileExpression) String() string {
	return fmt.Sprintf(
		"while %s do %s",
		we.If.String(),
		we.Loop.String(),
	)
}

type IndexExpression struct {
	Token lexer.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	return fmt.Sprintf(
		"%s[%s]",
		ie.Left.String(),
		ie.Index.String(),
	)
}

type SliceExpression struct {
	Token lexer.Token
	Left  Expression
	Start Expression
	End   Expression
}

func (se *SliceExpression) expressionNode() {}
func (se *SliceExpression) String() string {
	return fmt.Sprintf(
		"%s[%s:%s]",
		se.Left.String(),
		se.Start.String(),
		se.End.String(),
	)
}

type PropertyExpression struct {
	Token    lexer.Token
	Left     Expression
	Property *Identifier
}

func (pe *PropertyExpression) expressionNode() {}
func (pe *PropertyExpression) String() string {
	return fmt.Sprintf(
		"%s.%s",
		pe.Left.String(),
		pe.Property.String(),
	)
}

type KeyExpression struct {
	Token lexer.Token
	Left  Expression
	Key   Expression
}

func (ke *KeyExpression) expressionNode() {}
func (ke *KeyExpression) String() string {
	return fmt.Sprintf(
		"%s{%s}",
		ke.Left.String(),
		ke.Key.String(),
	)
}
