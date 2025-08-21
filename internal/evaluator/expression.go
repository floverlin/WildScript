package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/environment"
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

	"<":  "__lt",
	"<=": "__le",
	">":  "__gt",
	">=": "__ge",
	"==": "__eq",
	"!=": "__ne",
}

var unOps = map[string]string{
	"-":   "__unm",
	"not": "__not",
}

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) environment.Object {
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

	f, err := lookup(left, binOps[node.Operator])
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
) environment.Object {
	right := e.Eval(node.Right)

	f, err := lookup(right, unOps[node.Operator])
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
) environment.Object {
	left := e.Eval(node.Function)

	var self environment.Object
	if prop, ok := node.Function.(*ast.AttributeExpression); ok {
		self = e.Eval(prop.Left)
	} else if doc, ok := left.(*environment.Doc); ok {
		self = doc
	} else {
		self = environment.GLOBAL_NIL
	}

	args := e.evalExpressions(node.Arguments)
	var result environment.Object
	var err error

	if f, ok := left.(*environment.Func); ok {
		result, err = f.Call(e, self, args...)
	} else {
		f, metaErr := e.attribute(left, "__call")
		if metaErr != nil {
			lib.Die(
				node.Token,
				metaErr.Error(),
			)
		}
		result, err = f.(*environment.Func).Call(e, self, args...)
	}

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
) environment.Object {
	cond := e.Eval(node.If)

	if cond.Type() != environment.BOOL {
		lib.Die(
			node.Token,
			"non bool condition %s",
			cond.Type(),
		)
	}

	if cond.(*environment.Bool).Value {
		return e.Eval(node.Then)
	} else {
		return e.Eval(node.Else)
	}
}
