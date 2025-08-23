package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/environment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalAttributeExpression(
	node *ast.AttributeExpression,
) environment.Object {
	object := e.Eval(node.Left)
	prop := environment.NewString(node.Attribute.Value)

	result, err := environment.MetaCall(object, "__attribute", e, nil, prop)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}
	return result
}

func (e *Evaluator) evalIndexExpression(
	node *ast.IndexExpression,
) environment.Object {
	left := e.Eval(node.Left)

	index := e.Eval(node.Index)
	if index.Type() != environment.NUMBER &&
		index.Type() != environment.NIL {
		lib.Die(
			node.Token,
			"non num index",
		)
	}

	if index.Type() == environment.NIL {
		result, err := environment.MetaCall(left, "__list", e, nil)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}

		return result
	}

	result, err := environment.MetaCall(left, "__index", e, nil, index)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}

func (e *Evaluator) evalSliceExpression(
	node *ast.SliceExpression,
) environment.Object {
	left := e.Eval(node.Left)
	start := e.Eval(node.Start)
	end := e.Eval(node.End)

	if (start.Type() != environment.NUMBER &&
		start.Type() != environment.NIL) ||
		(end.Type() != environment.NUMBER &&
			end.Type() != environment.NIL) {
		lib.Die(
			node.Token,
			"non num index",
		)
	}

	result, err := environment.MetaCall(left, "__slice", e, nil, start, end)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}

func (e *Evaluator) evalKeyExpression(
	node *ast.KeyExpression,
) environment.Object {
	left := e.Eval(node.Left)
	key := e.Eval(node.Key)

	if key.Type() == environment.NIL {
		result, err := environment.MetaCall(left, "__dict", e, nil)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}

		return result
	}

	result, err := environment.MetaCall(left, "__key", e, nil, key)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}
