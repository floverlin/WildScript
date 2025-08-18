package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalWhileExpression(
	node *ast.WhileExpression,
) enviroment.Object {
	var result enviroment.Object
	condObject := e.Eval(node.If)
	cond, ok := condObject.(*enviroment.Bool)
	if !ok {
		lib.Die(
			node.Token,
			"non bool condition",
		)
	}
	for cond.Value {
		result = e.Eval(node.Loop)
		condObject = e.Eval(node.If)
		cond, ok = condObject.(*enviroment.Bool)
		if !ok {
			lib.Die(
				node.Token,
				"non bool condition",
			)
		}
	}

	return result
}
