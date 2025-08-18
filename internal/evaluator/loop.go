package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalWhileExpression(
	node *ast.WhileExpression,
) enviroment.Object {

	lib.Die(node.Token, "not exists")
	return nil
}
