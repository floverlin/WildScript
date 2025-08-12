package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/logger"
)

type Evaluator struct {
	env *enviroment.Enviroment
}

func New() *Evaluator {
	return &Evaluator{env: enviroment.New()}
}

func (e *Evaluator) Eval(node ast.Node) enviroment.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.BlockExpression:
		return e.evalBlockExpression(node)

	case *ast.AssignStatement:
		return e.evalAssignStatement(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)

	case *ast.InfixExpression:
		return e.evalInfixExpression(node)
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node)
	case *ast.CallExpression:
		return e.evalCallExpression(node)

	case *ast.Identifier:
		return e.evalIdentifier(node)

	case *ast.FloatLiteral:
		return &enviroment.Num{Value: node.Value}
	case *ast.StringLiteral:
		return &enviroment.Str{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value {
			return &e.env.Single().True
		} else {
			return &e.env.Single().False
		}
	case *ast.NilLiteral:
		return &e.env.Single().Nil
	default:
		panic("unknown node type")
	}
}

func (e *Evaluator) evalProgram(program *ast.Program) enviroment.Object {
	var result enviroment.Object
	for _, stmt := range program.Statements {
		result = e.Eval(stmt)
	}
	return result
}

func (e *Evaluator) evalBlockExpression(
	block *ast.BlockExpression,
) enviroment.Object {
	outerEnv := e.env
	e.env = enviroment.NewBlockEnviroment(outerEnv)

	var result enviroment.Object
	for _, stmt := range block.Statements {
		result = e.Eval(stmt)
	}

	e.env = outerEnv

	return result
}

func (e *Evaluator) evalExpressions(
	exprs []ast.Expression,
) []enviroment.Object {
	var result []enviroment.Object
	for _, expr := range exprs {
		evaluated := e.Eval(expr)
		result = append(result, evaluated)
	}
	return result
}

func (e *Evaluator) evalIdentifier(
	identifier *ast.Identifier,
) enviroment.Object {
	var val enviroment.Object
	var ok bool

	if identifier.IsOuter {
		val, ok = e.env.GetOuter(identifier.Value)
	} else {
		val, ok = e.env.Get(identifier.Value)
	}

	if ok {
		return val
	}
	panic(
		logger.Slog(
			identifier.Token.Line,
			identifier.Token.Column,
			"undefined variable: %s",
			identifier.Value,
		),
	)
}

func (e *Evaluator) evalAssignStatement(
	stmt *ast.AssignStatement,
) enviroment.Object {
	ident, typeOk := stmt.Left.(*ast.Identifier)
	if !typeOk {
		panic(
			logger.Slog(
				stmt.Token.Line,
				stmt.Token.Column,
				"expected identifier",
			),
		)
	}

	right := e.Eval(stmt.Right)

	var ok bool
	var result enviroment.Object

	if ident.IsOuter {
		result, ok = e.env.SetOuter(ident.Value, right)
	} else {
		result = e.env.Set(ident.Value, right)
		ok = true
	}

	if !ok {
		panic(
			logger.Slog(
				ident.Token.Line,
				ident.Token.Column,
				"undefined variable: %s",
				ident.Value,
			),
		)
	}

	return result
}
