package evaluator

import (
	"errors"
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func lookupMeta(
	object enviroment.Object,
	metaName string,
) (enviroment.MetaFunc, error) {
	meta, ok := enviroment.DefaultMeta[object.Type()]
	if !ok {
		return nil, errors.New("unsupported type")
	}
	f, ok := meta[metaName]
	if !ok {
		return nil, errors.New("unsupported operation")
	}
	return f, nil
}

func (e *Evaluator) evalPropertyExpression(
	node *ast.PropertyExpression,
) enviroment.Object {
	object := e.Eval(node.Left)
	prop := &enviroment.Str{Value: node.Property.Value}

	f, err := lookupMeta(object, "__property")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f(object, prop)
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
) enviroment.Object {
	left := e.Eval(node.Left)

	index := e.Eval(node.Index)
	if index.Type() != enviroment.NUM {
		lib.Die(
			node.Token,
			"non num index",
		)
	}

	f, err := lookupMeta(left, "__index")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f(left, index)
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

	f, err := lookupMeta(left, "__slice")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f(left, start, end)
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
) enviroment.Object {
	left := e.Eval(node.Left)
	key := e.Eval(node.Key)

	f, err := lookupMeta(left, "__key")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f(left, key)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}
