package ast

import (
	"fmt"
	"strconv"
	"strings"
	"arc/internal/lexer"
)

type Identifier struct {
	Token   lexer.Token
	Value   string
	IsOuter bool
	IsRune  bool
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	if i.IsOuter {
		return "&" + i.Value
	} else if i.IsRune {
		return "@" + i.Value
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

type ListLiteral struct {
	Token    lexer.Token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode() {}
func (ll *ListLiteral) String() string {
	var out strings.Builder
	out.WriteByte('[')
	for idx, elem := range ll.Elements {
		out.WriteString(elem.String())
		if idx != len(ll.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteByte(']')
	return out.String()
}

type ObjectField struct {
	Key   *Identifier
	Value Expression
}

type ObjectLiteral struct {
	Token  lexer.Token
	Fields []*ObjectField
}

func (ol *ObjectLiteral) expressionNode() {}
func (ol *ObjectLiteral) String() string {
	var out strings.Builder
	out.WriteString("{ ")
	for idx, field := range ol.Fields {
		if idx != 0 {
			out.WriteString("    ")
		}
		out.WriteString(
			fmt.Sprintf(
				"%s: %s",
				field.Key.String(),
				field.Value.String(),
			),
		)
		if idx != len(ol.Fields)-1 {
			out.WriteString(",\n")
		}
	}
	out.WriteString(" }")
	return out.String()
}
