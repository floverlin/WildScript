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
	for idx, stmt := range p.Statements {
		out.WriteString(stmt.String())
		if idx != len(p.Statements)-1 {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func joiner(args ...string) string {
	return strings.Join(args, " ")
}
