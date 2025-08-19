package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	right := e.Eval(node.Right)

	if left.Type() != right.Type() {
		lib.Die(
			node.Token,
			"wrong operands type %s and %s",
			left.Type(),
			right.Type(),
		)
	}

	switch left.Type() {
	case enviroment.NUM:
		return evalBinary(
			left.(*enviroment.Num).Value,
			right.(*enviroment.Num).Value,
			numOps,
			node,
		)
	case enviroment.BOOL:
		return evalBinary(
			left.(*enviroment.Bool).Value,
			right.(*enviroment.Bool).Value,
			boolOps,
			node,
		)
	case enviroment.STR:
		return evalBinary(
			left.(*enviroment.Str).Value,
			right.(*enviroment.Str).Value,
			strOps,
			node,
		)
	default:
		lib.Die(
			node.Token,
			"unsupported %s operands type %s",
			node.Operator,
			left.Type(),
		)
		return nil
	}
}

func (e *Evaluator) evalPrefixExpression(
	node *ast.PrefixExpression,
) enviroment.Object {
	right := e.Eval(node.Right)

	switch node.Operator {
	case "not":
		if boolObject, ok := right.(*enviroment.Bool); ok {
			if !boolObject.Value {
				return enviroment.GLOBAL_TRUE
			} else {
				return enviroment.GLOBAL_FALSE
			}
		} else {
			lib.Die(
				node.Token,
				"unsupperted not operand type %s",
				right.Type(),
			)
			return nil
		}
	case "-":
		if numObject, ok := right.(*enviroment.Num); ok {
			return &enviroment.Num{Value: -numObject.Value}
		} else {
			lib.Die(
				node.Token,
				"unsupperted minus operand type %s",
				right.Type(),
			)
			return nil
		}
	default:
		lib.Die(
			node.Token,
			"unknown prefix operator %s",
			node.Operator,
		)
		return nil
	}
}

func (e *Evaluator) evalCallExpression(
	node *ast.CallExpression,
) enviroment.Object {
	callable := e.Eval(node.Function)

	if callable.Type() != enviroment.FUNCTION &&
		callable.Type() != enviroment.NATIVE_FUNCTION {
		lib.Die(
			node.Token,
			"callable %s is not a function",
			callable.Inspect(),
		)
	}

	args := e.evalExpressions(node.Arguments)

	if f, ok := callable.(*enviroment.NativeFunction); ok {
		return f.Native(e, args...)
	}

	f := callable.(*enviroment.Function)

	if len(args) != len(f.Parameters) {
		lib.Die(
			node.Token,
			"function want %d argument(s) got %d",
			len(args),
			len(f.Parameters),
		)
	}

	fArgs := map[string]enviroment.Object{} // args
	for idx, arg := range args {
		fArgs[f.Parameters[idx].Value] = arg
	}

	result := e.EvalBlock(f.Body, f.Enviroment, fArgs)

	if result.Type() == enviroment.SIGNAL {
		if ret, ok := result.(*enviroment.Return); ok {
			result = ret.Value
		} else {
			lib.Die(
				node.Token,
				"continue or break in function",
			)
		}
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
