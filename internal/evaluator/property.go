package evaluator

import (
	"errors"
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func lookup(
	object enviroment.Object,
	metaName string,
) (enviroment.Callable, error) {
	meta, ok := enviroment.DefaultMeta[object.Type()]
	if !ok {
		return nil, errors.New("unsupported type")
	}
	f, ok := meta[metaName]
	if !ok {
		return nil, errors.New("unsupported operation")
	}
	return &enviroment.Func{Native: f}, nil
}

func (e *Evaluator) attribute(object enviroment.Object, attr string) (enviroment.Object, error) {
	f, err := lookup(object, "__attribute")
	if err != nil {
		return nil, err
	}

	result, err := f.Call(e, object, &enviroment.Str{Value: attr})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Evaluator) evalAttributeExpression(
	node *ast.AttributeExpression,
) enviroment.Object {
	object := e.Eval(node.Left)
	prop := &enviroment.Str{Value: node.Attribute.Value}

	f, err := lookup(object, "__attribute")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, object, prop)
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

	f, err := lookup(left, "__index")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, left, index)
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

	if (start.Type() != enviroment.NUM &&
		start.Type() != enviroment.NIL) ||
		(end.Type() != enviroment.NUM &&
			end.Type() != enviroment.NIL) {
		lib.Die(
			node.Token,
			"non num index",
		)
	}

	f, err := lookup(left, "__slice")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, left, start, end)
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

	if key.Type() == enviroment.NIL {
		f, err := lookup(left, "__dict")
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}

		result, err := f.Call(e, left)
		if err != nil {
			lib.Die(
				node.Token,
				err.Error(),
			)
		}

		return result
	}

	f, err := lookup(left, "__key")
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, left, key)
	if err != nil {
		lib.Die(
			node.Token,
			err.Error(),
		)
	}

	return result
}
