package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/environment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalWhileStatement(
	node *ast.WhileStatement,
) environment.Object {
	cond, err := environment.CheckBool(e.Eval(node.If))
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	var iters float64
	for cond {
		e.Eval(node.Loop)
		iters++

		cond, err = environment.CheckBool(e.Eval(node.If))
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
	}

	return environment.NewNumber(iters)
}

func (e *Evaluator) evalRepeatStatement(
	node *ast.RepeatStatement,
) environment.Object {
	var iters float64
	for {
		e.Eval(node.Loop)
		iters++

		cond, err := environment.CheckBool(e.Eval(node.Until))
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
		if !cond {
			break
		}
	}

	return environment.NewNumber(iters)
}

func (e *Evaluator) evalForStatement(
	node *ast.ForStatement,
) environment.Object {
	iterable := e.Eval(node.Iterable)

	iter, err := environment.MetaCall(iterable, "__iter", e, iterable)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	var iters float64
	for {
		next, err := environment.MetaCall(iter, "__next", e, iter)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
		value, cont, err := environment.UnpackResult(next)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
		if !cont {
			break
		}
		args := map[string]environment.Object{}
		if node.Value != nil {
			args[node.Value.Value] = value
		}
		e.EvalBlock(node.Loop, e.env, args)
		iters++
	}

	return environment.NewNumber(iters)
}
