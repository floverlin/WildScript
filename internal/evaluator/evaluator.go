package evaluator

import (
	"fmt"
	"os"
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
	"wildscript/internal/parser"
)

type Arguments = map[string]enviroment.Object

type Evaluator struct {
	env *enviroment.Enviroment
}

func New(env *enviroment.Enviroment, args, runes Arguments) *Evaluator {
	e := &Evaluator{env: enviroment.New(env)}
	for key, val := range args {
		e.env.Set(key, val)
	}
	for key, val := range runes {
		e.env.SetRune(key, val)
	}
	return e
}

func (e *Evaluator) Eval(node ast.Node) enviroment.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.BlockExpression:
		return e.EvalBlock(node, nil, nil)

	case *ast.AssignStatement:
		return e.evalAssignStatement(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.UseStatement:
		return e.evalUseStatement(node)
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
			return enviroment.Global[enviroment.GLOBAL_TRUE]
		} else {
			return enviroment.Global[enviroment.GLOBAL_FALSE]
		}
	case *ast.NilLiteral:
		return enviroment.Global[enviroment.GLOBAL_NIL]
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

func (e *Evaluator) EvalBlock(
	block *ast.BlockExpression,
	args Arguments,
	runes Arguments,
) enviroment.Object {
	var result enviroment.Object

	blockEval := New(e.env, args, runes)

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

func (e *Evaluator) evalUseStatement(
	node *ast.UseStatement,
) enviroment.Object {
	input, err := os.ReadFile(node.Name.Value + ".ws")
	if err != nil {
		panic(fmt.Sprintf("read module error: %s", err))
	}

	l := lexer.New(input)
	c := lexer.NewCollector(l)
	p := parser.New(c)
	mod := p.ParseProgram()

	modEv := New(nil, nil, nil)

	result := modEv.Eval(mod)

	e.env.Set(node.Name.Value, result)

	return enviroment.Global[enviroment.GLOBAL_NIL]
}

func (e *Evaluator) evalIdentifier(
	identifier *ast.Identifier,
) enviroment.Object {
	var val enviroment.Object
	var ok bool

	if identifier.IsRune {
		val, ok = e.env.GetRuneOuter(identifier.Value)
	} else if identifier.IsOuter {
		val, ok = e.env.GetOuter(identifier.Value)
	} else {
		val, ok = e.env.GetOuter(identifier.Value)
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
		if field.Key.IsRune {
			newObj.Runes[field.Key.Value] = value
		} else {
			newObj.Fields[field.Key.Value] = value
		}
	}

	return newObj
}
