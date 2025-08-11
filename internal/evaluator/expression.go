package evaluator

import (
	"math"
	"reflect"
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
		return e.evalFloatInfixExpression(left, right, node)
	case enviroment.BOOL_TYPE:
		return e.evalBooleanInfixExpression(left, right, node)
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
				return e.env.Single.True
			} else {
				return e.env.Single.False
			}
		case *enviroment.Str:
			if v.Value == "" {
				return e.env.Single.True
			} else {
				return e.env.Single.False
			}
		case *enviroment.Bool:
			if !v.Value {
				return e.env.Single.True
			} else {
				return e.env.Single.False
			}
		case *enviroment.Func:
			args := reflect.ValueOf(v.Fn).Type().NumIn()
			if args == 0 {
				return e.env.Single.True
			} else {
				return e.env.Single.False
			}
		case *enviroment.Nil:
			return e.env.Single.True
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
	function := e.Eval(node.Function)

	if function.Type() != enviroment.FUNC_TYPE {
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"not a function: %s",
				function.Type(),
			),
		)
	}

	args := e.evalExpressions(node.Arguments)

	return function.(*enviroment.Func).Fn(args...)
}

func (e *Evaluator) evalFloatInfixExpression(
	left, right enviroment.Object, 
	node *ast.InfixExpression) enviroment.Object {
	leftVal := left.(*enviroment.Num).Value
	rightVal := right.(*enviroment.Num).Value

	switch node.Operator {
	case "+":
		return &enviroment.Num{Value: leftVal + rightVal}
	case "-":
		return &enviroment.Num{Value: leftVal - rightVal}
	case "*":
		return &enviroment.Num{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			panic(
				logger.Slog(
					node.Token.Line,
					node.Token.Column,
					"division by zero",
				),
			)
		}
		return &enviroment.Num{Value: leftVal / rightVal}
	case "//":
		if rightVal == 0 {
			panic(
				logger.Slog(
					node.Token.Line,
					node.Token.Column,
					"division by zero",
				),
			)
		}
		return &enviroment.Num{Value: math.Floor(leftVal / rightVal)}
	case "%":
		if rightVal == 0 {
			panic(
				logger.Slog(
					node.Token.Line,
					node.Token.Column,
					"modulo by zero",
				),
			)
		}
		return &enviroment.Num{Value: math.Mod(leftVal, rightVal)}
	case "^":
		return &enviroment.Num{Value: math.Pow(leftVal, rightVal)}

	case "==":
		return &enviroment.Bool{Value: leftVal == rightVal}
	case "!=":
		return &enviroment.Bool{Value: leftVal != rightVal}
	case "<":
		return &enviroment.Bool{Value: leftVal < rightVal}
	case ">":
		return &enviroment.Bool{Value: leftVal > rightVal}
	case "<=":
		return &enviroment.Bool{Value: leftVal <= rightVal}
	case ">=":
		return &enviroment.Bool{Value: leftVal >= rightVal}
	default:
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"unsupported operator: %s",
				node.Operator,
			),
		)
	}
}

func (e *Evaluator) evalBooleanInfixExpression(left, right enviroment.Object, 
	node *ast.InfixExpression,
	) enviroment.Object {
	leftVal := left.(*enviroment.Bool).Value
	rightVal := right.(*enviroment.Bool).Value

	switch node.Operator {
	case "==":
		return &enviroment.Bool{Value: leftVal == rightVal}
	case "!=":
		return &enviroment.Bool{Value: leftVal != rightVal}
	case "||":
		return &enviroment.Bool{Value: leftVal || rightVal}
	case "&&":
		return &enviroment.Bool{Value: leftVal && rightVal}
	default:
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"unsupported operator: %s",
				node.Operator,
			),
		)
	}
}
