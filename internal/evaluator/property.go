package evaluator

import (
	"arc/internal/ast"
	"arc/internal/enviroment"
	"arc/internal/logger"
)

func (e *Evaluator) evalPropertyAccessExpression(
	node *ast.PropertyAccessExpression,
) enviroment.Object {
	object := e.Eval(node.Object, nil)
	propIdent := node.Property.Value

	// find field in obj
	if object.Type() == enviroment.OBJ_TYPE {
		obj := object.(*enviroment.Obj)
		if prop, ok := obj.Fields[propIdent]; ok {
			if prop.Type() == enviroment.FUNC_TYPE {
				self := enviroment.NewRune("self")
				self.Set(obj)
			}
			return prop
		}
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
	left := e.Eval(node.Left, nil)
	index := e.Eval(node.Index, nil)

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
	left := e.Eval(node.Left, nil)
	start := e.Eval(node.Start, nil)
	end := e.Eval(node.End, nil)

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
