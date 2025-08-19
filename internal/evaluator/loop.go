package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalWhileStatement(
	node *ast.WhileStatement,
) enviroment.Object {
	condObject := e.Eval(node.If)
	cond, ok := condObject.(*enviroment.Bool)
	if !ok {
		lib.Die(
			node.Token,
			"non bool condition",
		)
	}
	var iters float64
	for cond.Value {
		e.Eval(node.Loop)
		iters++

		condObject = e.Eval(node.If)
		cond, ok = condObject.(*enviroment.Bool)
		if !ok {
			lib.Die(
				node.Token,
				"non bool condition",
			)
		}
	}

	return &enviroment.Num{Value: iters}
}

func (e *Evaluator) evalRepeatStatement(
	node *ast.RepeatStatement,
) enviroment.Object {
	var iters float64
	for {
		e.Eval(node.Loop)
		iters++

		condObject := e.Eval(node.Until)
		cond, ok := condObject.(*enviroment.Bool)
		if !ok {
			lib.Die(
				node.Token,
				"non bool condition",
			)
		}
		if !cond.Value {
			break
		}
	}

	return &enviroment.Num{Value: iters}
}
