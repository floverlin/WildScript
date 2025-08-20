package evaluator

import (
	"fmt"
	"os"
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lexer"
	"wildscript/internal/lib"
	"wildscript/internal/parser"
)

type Evaluator struct {
	env *enviroment.Enviroment
}

func New(env *enviroment.Enviroment) *Evaluator {
	e := &Evaluator{env: enviroment.New(env)}
	return e
}

func (e *Evaluator) Eval(node ast.Node) enviroment.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.BlockExpression:
		return e.EvalBlock(node, e.env, nil)

	case *ast.AssignStatement:
		return e.evalAssignStatement(node)
	case *ast.LetStatement:
		return e.evalLetStatement(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.ImportStatement:
		return e.evalImportStatement(node)
	case *ast.FunctionStatement:
		return e.evalLetStatement(
			&ast.LetStatement{
				Token: node.Token,
				Left:  node.Identifier,
				Right: node.Function,
			},
		)
	case *ast.ReturnStatement:
		return &enviroment.Return{Value: e.Eval(node.Value)}
	case *ast.ContinueStatement:
		return &enviroment.Continue{}
	case *ast.BreakStatement:
		return &enviroment.Break{}
	case *ast.WhileStatement:
		return e.evalWhileStatement(node)
	case *ast.RepeatStatement:
		return e.evalRepeatStatement(node)

	case *ast.InfixExpression:
		return e.evalInfixExpression(node)
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node)
	case *ast.IfExpression:
		return e.evalIfExpression(node)
	case *ast.CallExpression:
		return e.evalCallExpression(node)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node)
	case *ast.SliceExpression:
		return e.evalSliceExpression(node)
	case *ast.AttributeExpression:
		return e.evalAttributeExpression(node)
	case *ast.KeyExpression:
		return e.evalKeyExpression(node)

	case *ast.Identifier:
		return e.evalIdentifier(node)

	case *ast.DocumentLiteral:
		return e.evalDocumentLiteral(node)

	case *ast.FunctionLiteral:
		return &enviroment.Func{
			Parameters: node.Parameters,
			Body:       node.Body,
			Enviroment: e.env,
			Impl:       node.Impl,
		}
	case *ast.NumberLiteral:
		return &enviroment.Num{Value: node.Value}
	case *ast.StringLiteral:
		return &enviroment.Str{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value {
			return enviroment.GLOBAL_TRUE
		} else {
			return enviroment.GLOBAL_FALSE
		}
	case *ast.NilLiteral:
		return enviroment.GLOBAL_NIL
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
	outer *enviroment.Enviroment,
	args map[string]enviroment.Object,
) enviroment.Object {
	var result enviroment.Object

	blockEval := New(outer)
	for key, val := range args {
		blockEval.env.Create(key, val)
	}

	for _, stmt := range block.Statements {
		result = blockEval.Eval(stmt)

		if result.Type() == enviroment.SIGNAL {
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
		result = append(result, e.Eval(expr))
	}
	return result
}

func (e *Evaluator) evalImportStatement(
	node *ast.ImportStatement,
) enviroment.Object {
	input, err := os.ReadFile(node.Module.Value + lib.EXT)
	if err != nil {
		panic(fmt.Sprintf("read module error: %s", err))
	}

	l := lexer.New(input)
	p := parser.New(l)
	mod := p.ParseProgram()

	modEv := New(nil)

	result := modEv.Eval(mod)

	e.env.Set(node.Module.Value, result)

	return enviroment.GLOBAL_NIL
}

func (e *Evaluator) evalIdentifier(
	identifier *ast.Identifier,
) enviroment.Object {
	val, ok := e.env.Get(identifier.Value)

	if ok {
		return val
	}
	lib.Die(
		identifier.Token,
		"undefined variable: %s",
		identifier.Value,
	)
	return nil
}

func (e *Evaluator) evalDocumentLiteral(
	node *ast.DocumentLiteral,
) enviroment.Object {
	doc := enviroment.NewDoc()
	for _, elem := range node.Elements {
		switch elem.Type {
		case ast.LIST:
			val := e.Eval(elem.Value)
			doc.List = append(doc.List, val)
		case ast.DICT:
			key, val := e.Eval(elem.Key), e.Eval(elem.Value)
			doc.Dict.Set(key, val)
		case ast.PROP:
			key := elem.Key.(*ast.Identifier).Value
			val := e.Eval(elem.Value)
			doc.Attrs[key] = val
		}
	}
	return doc
}
