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

func (e *Evaluator) evalForStatement(
	node *ast.ForStatement,
) enviroment.Object {
	iterable := e.Eval(node.Iterable)

	f, ok := enviroment.LookupMeta(iterable, "__iter")
	if !ok {
		lib.Die(
			node.Token,
			"not iterable",
		)
	}

	iter, err := f.(*enviroment.Func).Call(e, iterable)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	var iters float64
	for {
		result, err := iter.(*enviroment.Doc).
			Attrs["__next"].(*enviroment.Func).
			Native(e, iter)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
		ok := result.(*enviroment.Doc).Attrs["ok"]
		if !ok.(*enviroment.Bool).Value {
			break
		}
		args := map[string]enviroment.Object{}
		if node.Value != nil {
			value := result.(*enviroment.Doc).Attrs["value"]
			args[node.Value.Value] = value
		}
		e.EvalBlock(node.Loop, e.env, args)
		iters++
	}

	return &enviroment.Num{Value: iters}
}
