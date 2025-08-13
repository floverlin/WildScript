package ast

import (
	"fmt"
	"strconv"
	"strings"
	"wildscript/internal/lexer"
)

type Identifier struct {
	Token   lexer.Token
	Value   string
	IsOuter bool
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	if i.IsOuter {
		return "&" + i.Value
	}
	return i.Value
}

type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string {
	return strconv.FormatFloat(fl.Value, 'g', -1, 64)
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Value)
}

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	return strconv.FormatBool(bl.Value)
}

type NilLiteral struct {
	Token lexer.Token
}

func (nl *NilLiteral) expressionNode() {}
func (nl *NilLiteral) String() string {
	return "nil"
}

type FuncLiteral struct {
	Token      lexer.Token
	Parameters []*Identifier
	Body       *BlockExpression
}

func (fl *FuncLiteral) expressionNode() {}
func (fl *FuncLiteral) String() string {
	var out strings.Builder
	out.WriteString("fn(")
	for idx, param := range fl.Parameters {
		out.WriteString(param.String())
		if idx != len(fl.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") " + fl.Body.String())
	return out.String()
}
