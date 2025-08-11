package ast

import (
	"fmt"
	"strconv"
	"wildscript/internal/lexer"
)

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
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
