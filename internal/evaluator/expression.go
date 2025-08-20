package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

var binOps = map[string]string{
	"+":  "__add",
	"-":  "__sub",
	"*":  "__mul",
	"/":  "__div",
	"//": "__floor_div",
	"%":  "__mod",
	"^":  "__pow",
}

var unOps = map[string]string{
	"-":   "__unm",
	"not": "__not",
}

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	right := e.Eval(node.Right)

	if left.Type() != right.Type() {
		lib.Die(
			node.Token,
			"non equal operands type %s and %s",
			left.Type(),
			right.Type(),
		)
	}

	f, err := lookupMeta(left, binOps[node.Operator])
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, left, right)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}

func (e *Evaluator) evalPrefixExpression(
	node *ast.PrefixExpression,
) enviroment.Object {
	right := e.Eval(node.Right)

	f, err := lookupMeta(right, unOps[node.Operator])
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, right)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}

func (e *Evaluator) evalCallExpression(
	node *ast.CallExpression,
) enviroment.Object {
	var self enviroment.Object
	if prop, ok := node.Function.(*ast.PropertyExpression); ok {
		self = e.Eval(prop.Left)
	} else {
		self = enviroment.GLOBAL_NIL
	}

	var t any
	left := e.Eval(node.Function)

	if left.Type() == enviroment.DOC {
		f, err := lookupMeta(left, "__call")
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}
		t = f
	} else {
		t = left
	}

	f, ok := t.(*enviroment.Func)
	if !ok {
		lib.Die(
			node.Token,
			"callable %s is not a function",
			left.Inspect(),
		)
	}

	args := e.evalExpressions(node.Arguments)

	if f.Impl == ast.NATIVE {
		return f.Native(args...)
	}

	if f.Impl == ast.METHOD {
		args = append([]enviroment.Object{self}, args...)
	}

	result, err := f.Call(e, args...)

	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}

func (e *Evaluator) evalIfExpression(
	node *ast.IfExpression,
) enviroment.Object {
	cond := e.Eval(node.If)

	if cond.Type() != enviroment.BOOL {
		lib.Die(
			node.Token,
			"non bool condition %s",
			cond.Type(),
		)
	}

	if cond.(*enviroment.Bool).Value {
		return e.Eval(node.Then)
	} else {
		return e.Eval(node.Else)
	}
}
