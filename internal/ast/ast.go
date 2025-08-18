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
	var sb strings.Builder
	for idx, stmt := range p.Statements {
		sb.WriteString(stmt.String())
		if idx != len(p.Statements)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
