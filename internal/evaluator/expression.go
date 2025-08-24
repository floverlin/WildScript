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

	if node.Operator == "and" ||
		node.Operator == "or" {
		if left.Type() != environment.BOOLEAN {
			lib.Die(
				node.Token,
				"non boolean condition",
			)
		}
		if node.Operator == "and" {
			if b, _ := environment.CheckBool(left); b {
				return right
			} else {
				return left
			}
		}
		if node.Operator == "or" {
			if b, _ := environment.CheckBool(left); b {
				return left
			} else {
				return right
			}
		}
	}

	result, err := environment.MetaCall(left, binOps[node.Operator], e, nil, right)
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

	result, err := environment.MetaCall(right, unOps[node.Operator], e, nil)
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
	} else {
		self = nil
	}

	args := e.evalExpressions(node.Arguments)

	result, err := environment.MetaCall(left, "__call", e, self, args...)

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
	cond, err := environment.CheckBool(e.Eval(node.If))
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	if cond {
		return e.Eval(node.Then)
	} else {
		return e.Eval(node.Else)
	}
}
