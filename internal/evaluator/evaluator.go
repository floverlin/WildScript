package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/logger"
)

type Arguments = map[string]enviroment.Object

type Evaluator struct {
	env *enviroment.Enviroment
}

func New(env *enviroment.Enviroment, args Arguments) *Evaluator {
	e := &Evaluator{env: enviroment.New(env)}
	for key, val := range args {
		e.env.Set(key, val)
	}
	return e
}

func (e *Evaluator) Eval(node ast.Node) enviroment.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.BlockExpression:
		return e.EvalBlock(node, nil)

	case *ast.AssignStatement:
		return e.evalAssignStatement(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.FuncStatement:
		return e.evalAssignStatement(
			&ast.AssignStatement{
				Token: node.Token,
				Left:  node.Identifier,
				Right: node.Function,
			},
		)
	case *ast.ReturnStatement:
		return &enviroment.Return{Value: e.Eval(node.Value)}
	case *ast.ContinueStatement:
		return &enviroment.Continue{}

	case *ast.InfixExpression:
		return e.evalInfixExpression(node)
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node)
	case *ast.CallExpression:
		return e.evalCallExpression(node)
	case *ast.ConditionExpression:
		return e.evalConditionExpression(node)
	case *ast.ForExpression:
		return e.evalForExpression(node)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node)
	case *ast.SliceExpression:
		return e.evalSliceExpression(node)
	case *ast.PropertyAccessExpression:
		return e.evalPropertyAccessExpression(node)

	case *ast.Identifier:
		return e.evalIdentifier(node)

	case *ast.FuncLiteral:
		return &enviroment.Func{
			Parameters: node.Parameters,
			Body:       node.Body,
			Enviroment: e.env,
		}
	case *ast.ListLiteral:
		return e.evalListLiteral(node)
	case *ast.ObjectLiteral:
		return e.evalObjectLiteral(node)
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
		panic("unknown node type\n")
	}
}

func (e *Evaluator) evalProgram(program *ast.Program) enviroment.Object {
	var result enviroment.Object
	for _, stmt := range program.Statements {
		result = e.Eval(stmt)
	}
	return result
}

func (e *Evaluator) EvalBlock(
	block *ast.BlockExpression,
	args Arguments,
) enviroment.Object {
	var result enviroment.Object

	blockEval := New(e.env, args)

	for _, stmt := range block.Statements {
		result = blockEval.Eval(stmt)

		if result.Type() == enviroment.CONTROL_TYPE {
			break
		}
	}

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

	if identifier.IsRune {
		r, runeOk := enviroment.FindRune(identifier.Value)
		val = r.Get()
		ok = runeOk
	} else if identifier.IsOuter {
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

func (e *Evaluator) evalListLiteral(
	node *ast.ListLiteral,
) *enviroment.List {
	elems := []enviroment.Object{}

	for _, nodeElem := range node.Elements {
		elems = append(elems, e.Eval(nodeElem))
	}

	return &enviroment.List{Elements: elems}
}

func (e *Evaluator) evalObjectLiteral(
	node *ast.ObjectLiteral,
) enviroment.Object {
	newObj := enviroment.NewObj()

	for _, field := range node.Fields {
		value := e.Eval(field.Value)
		newObj.Fields[field.Key.Value] = value
	}

	return newObj
}
