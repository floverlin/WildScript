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
		return e.evalBlockExpression(node, nil)

	case *ast.AssignStatement:
		return e.evalAssignStatement(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
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

func (e *Evaluator) evalBlockExpression(
	block *ast.BlockExpression,
	args []blockArgument,
) enviroment.Object {
	outerEnv := e.env
	e.env = enviroment.NewBlockEnviroment(outerEnv)

	for _, arg := range args {
		e.env.Set(arg.Name, arg.Value)
	}

	var result enviroment.Object
	for _, stmt := range block.Statements {
		result = e.Eval(stmt)

		if result.Type() == enviroment.CONTROL_TYPE {
			break
		}
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
	fields := map[string]enviroment.Object{}

	for _, field := range node.Fields {
		value := e.Eval(field.Value)
		fields[field.Key.Value] = value
	}

	return &enviroment.Obj{Fields: fields}
}

func (e *Evaluator) evalPropertyAccessExpression(
	node *ast.PropertyAccessExpression,
) enviroment.Object {
	obj := e.Eval(node.Object)
	propIdent := node.Property.Value

	if obj.Type() == enviroment.OBJ_TYPE {
		obj := obj.(*enviroment.Obj)
		if prop, ok := obj.Fields[propIdent]; ok {
			return prop
		}
	}

	method := enviroment.FindMethod(obj, propIdent)
	// if not built-in set self
	if method == nil {
		panic("TODO METHOD NOT FOUND")
	}

	return method
}

func (e *Evaluator) evalIndexExpression(
	node *ast.IndexExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	index := e.Eval(node.Index)

	if index.Type() != enviroment.NUM_TYPE {
		panic("TODO")
	}
	idx := int(index.(*enviroment.Num).Value)

	var result enviroment.Object

	switch v := left.(type) {
	case *enviroment.Str:
		sl := []rune(v.Value)
		symbol := sl[idx]
		result = &enviroment.Str{
			Value: string(symbol),
		}
	case *enviroment.List:
		result = v.Elements[idx]
	default:
		panic("TODO")
	}

	return result
}

func (e *Evaluator) evalSliceExpression(
	node *ast.SliceExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	start := e.Eval(node.Start)
	end := e.Eval(node.End)

	if start.Type() != enviroment.NUM_TYPE ||
		end.Type() != enviroment.NUM_TYPE {
		panic("TODO")
	}
	startVal := int(start.(*enviroment.Num).Value)
	endVal := int(end.(*enviroment.Num).Value)

	var result enviroment.Object

	switch v := left.(type) {
	case *enviroment.Str:
		sl := []rune(v.Value)
		symbols := sl[startVal:endVal]
		result = &enviroment.Str{
			Value: string(symbols),
		}
	case *enviroment.List:
		return &enviroment.List{
			Elements: v.Elements[startVal:endVal],
		}
	default:
		panic("TODO")
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

	if ident.IsRune {
		r, runeOk := enviroment.FindRune(ident.Value)
		result = r.Set(right)
		ok = runeOk
	} else if ident.IsOuter {
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
