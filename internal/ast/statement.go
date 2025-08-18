package ast

import (
	"fmt"
	"wildscript/internal/lexer"
)

type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type AssignStatement struct {
	Token lexer.Token
	Left  Expression
	Right Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) String() string {
	return fmt.Sprintf("%s = %s", as.Left.String(), as.Right.String())
}

type FunctionStatement struct {
	Token      lexer.Token
	Identifier *Identifier
	Function   *FunctionLiteral
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) String() string {
	return fmt.Sprintf(
		"function %s%s",
		fs.Identifier.String(),
		fs.Function.String(),
	)
}

type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", rs.Value.String())
}

type ContinueStatement struct {
	Token lexer.Token
}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) String() string {
	return "continue"
}

type BreakStatement struct {
	Token lexer.Token
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return "break"
}

type ImportStatement struct {
	Token  lexer.Token
	Module *Identifier
}

func (is *ImportStatement) statementNode() {}
func (is *ImportStatement) String() string {
	return fmt.Sprintf("import %s", is.Module.Value)
}

type ExportStatement struct {
	Token lexer.Token
	Value Expression
}

func (es *ExportStatement) statementNode() {}
func (es *ExportStatement) String() string {
	return fmt.Sprintf("export %s", es.Value.String())
}
