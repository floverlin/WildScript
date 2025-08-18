package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func findDocumentPropertry(
	document *enviroment.Doc,
	prop string,
) enviroment.Object {
	if value, ok := document.Elements[prop]; ok {
		return value
	}
	return nil
}

func (e *Evaluator) evalPropertyExpression(
	node *ast.PropertyExpression,
) enviroment.Object {
	object := e.Eval(node.Left)
	prop := node.Property.Value

	// find elem in doc
	elem := findDocumentPropertry(
		object.(*enviroment.Doc),
		prop,
	)

	if elem != nil {
		return elem
	}

	lib.Die(
		node.Token,
		"property %s not exists in %s",
		prop,
		object.Type(),
	)
	return nil
}

func (e *Evaluator) evalIndexExpression(
	node *ast.IndexExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	index := e.Eval(node.Index)
	var idx int

	if index, ok := index.(*enviroment.Num); ok {
		idx = int(index.Value)
	} else {
		lib.Die(
			node.Token,
			"non num index",
		)
	}

	var result enviroment.Object

	switch object := left.(type) {
	case *enviroment.Str:
		sl := []rune(object.Value)
		symbol := sl[idx]
		result = &enviroment.Str{
			Value: string(symbol),
		}
	case *enviroment.Doc:
		result = object.List[idx]
	default:
		lib.Die(
			node.Token,
			"unsupported index access for %s",
			left.Type(),
		)
	}

	return result
}

func (e *Evaluator) evalSliceExpression(
	node *ast.SliceExpression,
) enviroment.Object {
	left := e.Eval(node.Left)
	start := e.Eval(node.Start)
	end := e.Eval(node.End)

	if start.Type() != enviroment.NUM ||
		end.Type() != enviroment.NUM {
		lib.Die(
			node.Token,
			"non num index",
		)
	}
	startIdx := int(start.(*enviroment.Num).Value)
	endIdx := int(end.(*enviroment.Num).Value)

	var result enviroment.Object

	switch object := left.(type) {
	case *enviroment.Str:
		sl := []rune(object.Value)
		symbols := sl[startIdx:endIdx]
		result = &enviroment.Str{
			Value: string(symbols),
		}
	case *enviroment.Doc:
		return &enviroment.Doc{
			List: object.List[startIdx:endIdx],
		}
	default:
		lib.Die(
			node.Token,
			"unsupported slice access for %s",
			left.Type(),
		)
	}

	return result
}
