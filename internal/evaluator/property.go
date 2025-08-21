package evaluator

import (
	"errors"
	"wildscript/internal/ast"
	"wildscript/internal/environment"
	"wildscript/internal/lib"
)

func lookup(
	object environment.Object,
	metaName string,
) (environment.Callable, error) {
	meta, ok := environment.DefaultMeta[object.Type()]
	if !ok {
		return nil, errors.New("unsupported type")
	}
	f, ok := meta[metaName]
	if !ok {
		return nil, errors.New("unsupported operation")
	}
	return &environment.Func{Native: f}, nil
}

func (e *Evaluator) attribute(
	object environment.Object,
	attr string,
) (environment.Object, error) {
	f, err := lookup(object, "__attribute")
	if err != nil {
		return nil, err
	}

	result, err := f.Call(e, object, &environment.Str{Value: attr})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Evaluator) evalAttributeExpression(
	node *ast.AttributeExpression,
) environment.Object {
	object := e.Eval(node.Left)
	prop := &environment.Str{Value: node.Attribute.Value}

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
) environment.Object {
	left := e.Eval(node.Left)

	index := e.Eval(node.Index)
	if index.Type() != environment.NUM {
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
) environment.Object {
	left := e.Eval(node.Left)
	start := e.Eval(node.Start)
	end := e.Eval(node.End)

	if (start.Type() != environment.NUM &&
		start.Type() != environment.NIL) ||
		(end.Type() != environment.NUM &&
			end.Type() != environment.NIL) {
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
) environment.Object {
	left := e.Eval(node.Left)
	key := e.Eval(node.Key)

	if key.Type() == environment.NIL {
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
