package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
)

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
