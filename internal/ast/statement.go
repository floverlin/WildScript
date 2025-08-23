package ast

import (
	"fmt"
	"strings"
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

type LetStatement struct {
	Token lexer.Token
	Left  *Identifier
	Right Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", ls.Left.String(), ls.Right.String())
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
	Module []*Identifier
}

func (is *ImportStatement) statementNode() {}
func (is *ImportStatement) String() string {
	var sb strings.Builder
	sb.WriteString("import ")
	for _, mod := range is.Module {
		sb.WriteString(mod.String() + ".")
	}
	result := sb.String()
	return result[:len(result)-1]
}

type ExportStatement struct {
	Token lexer.Token
	Value Expression
}

func (es *ExportStatement) statementNode() {}
func (es *ExportStatement) String() string {
	return fmt.Sprintf("export %s", es.Value.String())
}

type ForStatement struct {
	Token    lexer.Token
	Value    *Identifier
	Iterable Expression
	Loop     *BlockExpression
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	if fs.Value != nil {
		return fmt.Sprintf(
			"for %s in %s do %s",
			fs.Value.String(),
			fs.Iterable.String(),
			fs.Loop.String(),
		)
	} else {
		return fmt.Sprintf(
			"for %s do %s",
			fs.Iterable.String(),
			fs.Loop.String(),
		)
	}
}

type RepeatStatement struct {
	Token lexer.Token
	Until Expression
	Loop  *BlockExpression
}

func (rs *RepeatStatement) statementNode() {}
func (rs *RepeatStatement) String() string {
	return fmt.Sprintf(
		"repeat %s until %s",
		rs.Loop.String(),
		rs.Until.String(),
	)
}

type WhileStatement struct {
	Token lexer.Token
	If    Expression
	Loop  *BlockExpression
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf(
		"while %s do %s",
		ws.If.String(),
		ws.Loop.String(),
	)
}
