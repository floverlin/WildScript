package ast

import "strings"

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out strings.Builder
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteByte('\n')
	}
	return out.String()
}

func joiner(args ...string) string {
	return strings.Join(args, " ")
}
