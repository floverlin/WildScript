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
		return function.Builtin(args...)
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

	outerEnv := e.env            // save init env
	funcArgs := []funcArgument{} // env for args
	for idx, arg := range args {
		funcArgs = append(funcArgs, funcArgument{
			Name:  function.Parameters[idx].Value,
			Value: arg,
		})
	}

	e.env = function.Enviroment // closure

	result := e.evalBlockExpression(function.Body, funcArgs)

	e.env = outerEnv

	return result
}

type funcArgument struct {
	Name  string
	Value enviroment.Object
}

func (e *Evaluator) evalForExpression(
	node *ast.ForExpression,
) enviroment.Object {
	var result enviroment.Object = &e.env.Single().Nil

	cond := e.Eval(node.Condition)

	if cond.Type() == enviroment.BOOL_TYPE {
		for {
			result = e.Eval(node.Body)
			cond = e.Eval(node.Condition)
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

		for range iters {
			result = e.Eval(node.Body)
		}
	}

	return result
}

func evalBinary[T any](
	left, right T,
	ops map[string]func(T, T) (enviroment.Object, error),
	node *ast.InfixExpression,
) enviroment.Object {
	if f, ok := ops[node.Operator]; ok {
		obj, err := f(left, right)
		if err != nil {
			panic(
				logger.Slog(
					node.Token.Line,
					node.Token.Column,
					"%s: %s",
					err.Error(),
					node.Operator,
				),
			)
		}
		return obj
	}

	panic(
		logger.Slog(
			node.Token.Line,
			node.Token.Column,
			"unsupported operator: %s",
			node.Operator,
		),
	)
}
