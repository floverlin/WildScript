package ast

import (
	"fmt"
	"strconv"
	"strings"
	"wildscript/internal/lexer"
)

type ElementType int
type FunctionType string

const (
	LIST ElementType = iota
	PROP
	DICT
)

const (
	FUNCTION FunctionType = "function"
	LAMBDA   FunctionType = "lambda"
	METHOD   FunctionType = "method"
)

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	Token lexer.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string {
	return strconv.FormatFloat(nl.Value, 'g', -1, 64)
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

type FunctionLiteral struct {
	Token      lexer.Token
	Parameters []*Identifier
	Body       *BlockExpression
	Type       FunctionType
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var sb strings.Builder
	if fl.Type != FUNCTION {
		sb.WriteString(string(fl.Type))
	}
	sb.WriteString("(")
	for idx, param := range fl.Parameters {
		sb.WriteString(param.String())
		if idx != len(fl.Parameters)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") " + fl.Body.String())
	return sb.String()
}

type DocumentElement struct {
	Key   *Identifier
	Type  ElementType
	Value Expression
}

func (de *DocumentElement) String() string {
	switch de.Type {
	case PROP:
		return fmt.Sprintf(
			"%s = %s",
			de.Key.String(),
			de.Value.String(),
		)
	case LIST:
		return de.Value.String()
	case DICT:
		return fmt.Sprintf(
			"%s: %s",
			de.Key.String(),
			de.Value.String(),
		)
	default:
		return "error"
	}
}

type DocumentLiteral struct {
	Token    lexer.Token
	Elements []*DocumentElement
}

func (dl *DocumentLiteral) expressionNode() {}
func (dl *DocumentLiteral) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for idx, elem := range dl.Elements {
		sb.WriteString(elem.String())
		if idx != len(dl.Elements)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
