package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/logger"
)

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	right := e.Eval(node.Right)

	if left.Type() != right.Type() {
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"wrong operands type: %s %s",
				left.Type(),
				right.Type(),
			),
		)
	}

	switch left.Type() {
	case enviroment.NUM_TYPE:
		return evalBinary(
			left.(*enviroment.Num).Value,
			right.(*enviroment.Num).Value,
			numOps,
			node,
		)
	case enviroment.BOOL_TYPE:
		return evalBinary(
			left.(*enviroment.Bool).Value,
			right.(*enviroment.Bool).Value,
			boolOps,
			node,
		)
	case enviroment.STR_TYPE:
		return evalBinary(
			left.(*enviroment.Str).Value,
			right.(*enviroment.Str).Value,
			strOps,
			node,
		)
	default:
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"wrong operands type %s %s",
				left.Type(),
				right.Type(),
			),
		)
	}
}

func (e *Evaluator) evalPrefixExpression(
	node *ast.PrefixExpression,
) enviroment.Object {
	value := e.Eval(node.Right)

	switch node.Operator {
	case "!":
		switch v := value.(type) {
		case *enviroment.Num:
			if v.Value == 0 {
				return &e.env.Single().True
			} else {
				return &e.env.Single().False
			}
		case *enviroment.Str:
			if v.Value == "" {
				return &e.env.Single().True
			} else {
				return &e.env.Single().False
			}
		case *enviroment.Bool:
			if !v.Value {
				return &e.env.Single().True
			} else {
				return &e.env.Single().False
			}
		case *enviroment.Func:
			params := v.LenOfParameters()
			if params == 0 {
				return &e.env.Single().True
			} else {
				return &e.env.Single().False
			}
		case *enviroment.Nil:
			return &e.env.Single().True
		default:
			panic(
				logger.Slog(
					node.Token.Line,
					node.Token.Column,
					"unknown object type",
				),
			)
		}
	default:
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"unknown prefix operator",
			),
		)
	}
}

func (e *Evaluator) evalCallExpression(
	node *ast.CallExpression,
) enviroment.Object {
	callable := e.Eval(node.Function)

	if callable.Type() == enviroment.OBJ_TYPE {
		obj := callable.(*enviroment.Obj)
		call, ok := obj.Runes[enviroment.CALL_RUNE]
		if ok && call.Type() == enviroment.FUNC_TYPE {
			call.(*enviroment.Func).
				Enviroment.SetRune(enviroment.SELF_RUNE, obj)
			callable = call
		}
	}

	if callable.Type() != enviroment.FUNC_TYPE {
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"not callable %s",
				callable.Type(),
			),
		)
	}

	args := e.evalExpressions(node.Arguments)
	function := callable.(*enviroment.Func)

	if function.Builtin != nil {
		return function.Builtin(e, args...)
	}

	if len(args) != function.LenOfParameters() {
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"function want %d argument(s) got %d",
				len(args),
				function.LenOfParameters(),
			),
		)
	}

	funcArgs := Arguments{} // args
	for idx, arg := range args {
		funcArgs[function.Parameters[idx].Value] = arg
	}

	callEval := New(function.Enviroment, nil, nil) // closure

	result := callEval.EvalBlock(function.Body, funcArgs, nil)

	if result.Type() == enviroment.CONTROL_TYPE {
		if ret, ok := result.(*enviroment.Return); ok {
			result = ret.Value
		} else {
			panic(
				logger.Slog(
					node.Token.Line,
					node.Token.Column,
					"continue in non for block",
				),
			)
		}
	}

	return result
}

func (e *Evaluator) evalConditionExpression(
	node *ast.ConditionExpression,
) enviroment.Object {
	cond := e.Eval(node.Condition)

	if cond.Type() != enviroment.BOOL_TYPE {
		panic("TODO")
	}

	if cond.(*enviroment.Bool).Value {
		return e.Eval(node.Consequence)
	} else {
		return e.Eval(node.Alternative)
	}
}
