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

	meta, ok := enviroment.DefaultMeta[left.Type()]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported type",
		)
	}

	op := binOps[node.Operator]

	f, ok := meta[op]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported operation",
		)
	}

	result, err := f(left, right)
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

	meta, ok := enviroment.DefaultMeta[right.Type()]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported type",
		)
	}

	op := unOps[node.Operator]

	f, ok := meta[op]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported operation",
		)
	}

	result, err := f(right)
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

	callable := e.Eval(node.Function)

	f, ok := callable.(*enviroment.Func)
	if !ok {
		lib.Die(
			node.Token,
			"callable %s is not a function",
			callable.Inspect(),
		)
	}

	args := e.evalExpressions(node.Arguments)

	if f.Impl == ast.NATIVE {
		return f.Native(args...)
	}

	if f.Impl == ast.METHOD {
		args = append([]enviroment.Object{self}, args...)
	}

	if len(args) != len(f.Parameters) {
		lib.Die(
			node.Token,
			"function want %d argument(s) got %d",
			len(f.Parameters),
			len(args),
		)
	}

	fArgs := map[string]enviroment.Object{} // args
	for idx, arg := range args {
		fArgs[f.Parameters[idx].Value] = arg
	}

	result := e.EvalBlock(f.Body, f.Enviroment, fArgs)

	if result.Type() == enviroment.SIGNAL {
		if ret, ok := result.(*enviroment.Return); ok {
			return ret.Value
		} else {
			lib.Die(
				node.Token,
				"continue or break in function",
			)
		}
	}

	return enviroment.GLOBAL_NIL
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
