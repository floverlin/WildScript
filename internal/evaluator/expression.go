package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/logger"
)

func (e *Evaluator) evalInfixExpression(
	node *ast.InfixExpression,
) enviroment.Object {
	left := e.Eval(node.Left, nil)
	right := e.Eval(node.Right, nil)

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
	value := e.Eval(node.Right, nil)

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
	callable := e.Eval(node.Function, nil)

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

	outerEnv := e.env // save init env

	funcArgs := Arguments{} // args
	for idx, arg := range args {
		funcArgs[function.Parameters[idx].Value] = arg
	}

	e.env = function.Enviroment // closure

	result := e.evalBlockExpression(function.Body, funcArgs)

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

	e.env = outerEnv

	return result
}

func (e *Evaluator) evalForExpression(
	node *ast.ForExpression,
) enviroment.Object {
	var result enviroment.Object = &e.env.Single().Nil

	cond := e.Eval(node.Condition, nil)

	if cond.Type() == enviroment.BOOL_TYPE {
		for {
			result = e.Eval(node.Body, nil)
			cond = e.Eval(node.Condition, nil)
			if !cond.(*enviroment.Bool).Value {
				break
			}
		}

	} else {
		var iters int
		switch c := cond.(type) {
		case *enviroment.Num:
			iters = int(c.Value)
		case *enviroment.Str:
			iters = len(c.Value)
		case *enviroment.Nil:
			iters = 0
		case *enviroment.Func:
			iters = c.LenOfParameters()
		default:
			panic("TODO")
		}

		idxRune := enviroment.NewRune("idx")
		keyRune := enviroment.NewRune("key")
		valRune := enviroment.NewRune("val")

		for idx := range iters {
			idxRune.Set(&enviroment.Num{Value: float64(idx)})
			keyRune.Set(&enviroment.Num{Value: float64(idx)})
			valRune.Set(&enviroment.Num{Value: float64(idx)})

			result = e.evalBlockExpression(
				node.Body,
				nil,
			)
			if result.Type() == enviroment.CONTROL_TYPE {
				switch res := result.(type) {
				case *enviroment.Continue:
					continue
				case *enviroment.Return:
					result = res.Value
				}
				break
			}
		}
	}

	return result
}

func (e *Evaluator) evalConditionExpression(
	node *ast.ConditionExpression,
) enviroment.Object {
	cond := e.Eval(node.Condition, nil)

	if cond.Type() != enviroment.BOOL_TYPE {
		panic("TODO")
	}

	if cond.(*enviroment.Bool).Value {
		return e.Eval(node.Consequence, nil)
	} else {
		return e.Eval(node.Alternative, nil)
	}
}
