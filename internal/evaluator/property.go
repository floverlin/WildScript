package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/logger"
)

func findObjectPropertry(
	object enviroment.Object,
	prop *ast.Identifier,
) enviroment.Object {
	if object.Type() != enviroment.OBJ_TYPE {
		return nil
	}
	obj := object.(*enviroment.Obj)

	if prop.IsRune {
		if prop, ok := obj.Runes[prop.Value]; ok {
			return prop
		}
	} else {
		if prop, ok := obj.Fields[prop.Value]; ok {
			return prop
		}
	}

	if proto, ok := obj.Runes[enviroment.PROTO_RUNE]; ok {
		return findObjectPropertry(proto, prop)
	}

	return nil
}

func (e *Evaluator) evalPropertyAccessExpression(
	node *ast.PropertyAccessExpression,
) enviroment.Object {
	object := e.Eval(node.Object)
	propIdent := node.Property.Value

	if node.Property.IsOuter {
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"outer property not exists",
			),
		)
	}

	// find field in obj
	prop := findObjectPropertry(object, node.Property)
	if prop != nil {
		if prop.Type() == enviroment.FUNC_TYPE {
			prop := prop.(*enviroment.Func)
			prop.Enviroment.SetRune(enviroment.SELF_RUNE, object)
		}
		return prop
	}

	// find base type method
	method := enviroment.FindMethod(object, propIdent)

	if method == nil {
		panic(
			logger.Slog(
				node.Token.Line,
				node.Token.Column,
				"property %s not exists in %s",
				propIdent,
				object.Type(),
			),
		)
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
