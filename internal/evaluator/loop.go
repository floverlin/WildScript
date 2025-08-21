package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/environment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalWhileStatement(
	node *ast.WhileStatement,
) environment.Object {
	condObject := e.Eval(node.If)
	cond, ok := condObject.(*environment.Bool)
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
		cond, ok = condObject.(*environment.Bool)
		if !ok {
			lib.Die(
				node.Token,
				"non bool condition",
			)
		}
	}

	return &environment.Num{Value: iters}
}

func (e *Evaluator) evalRepeatStatement(
	node *ast.RepeatStatement,
) environment.Object {
	var iters float64
	for {
		e.Eval(node.Loop)
		iters++

		condObject := e.Eval(node.Until)
		cond, ok := condObject.(*environment.Bool)
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

	return &environment.Num{Value: iters}
}

func (e *Evaluator) evalForStatement(
	node *ast.ForStatement,
) environment.Object {
	iterable := e.Eval(node.Iterable)

	f, ok := environment.LookupMeta(iterable, "__iter")
	if !ok {
		lib.Die(
			node.Token,
			"not iterable",
		)
	}

	iter, err := f.(*environment.Func).Call(e, iterable)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	var iters float64
	for {
		result, err := iter.(*environment.Doc).
			Attrs["__next"].(*environment.Func).
			Native(e, iter)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
		ok := result.(*environment.Doc).Attrs["ok"]
		if !ok.(*environment.Bool).Value {
			break
		}
		args := map[string]environment.Object{}
		if node.Value != nil {
			value := result.(*environment.Doc).Attrs["value"]
			args[node.Value.Value] = value
		}
		e.EvalBlock(node.Loop, e.env, args)
		iters++
	}

	return &environment.Num{Value: iters}
}
