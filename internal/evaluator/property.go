package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalPropertyExpression(
	node *ast.PropertyExpression,
) enviroment.Object {
	object := e.Eval(node.Left)
	prop := &enviroment.Str{Value: node.Property.Value}

	meta, ok := enviroment.DefaultMeta[object.Type()]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported type",
		)
	}

	f, ok := meta["__property"]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported operation",
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

	meta, ok := enviroment.DefaultMeta[left.Type()]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported type",
		)
	}

	f, ok := meta["__index"]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported operation",
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

	meta, ok := enviroment.DefaultMeta[left.Type()]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported type",
		)
	}

	f, ok := meta["__slice"]
	if !ok {
		lib.Die(
			node.Token,
			"unsupported operation",
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

	doc, ok := left.(*enviroment.Doc)
	if !ok {
		lib.Die(
			node.Token,
			"key access to non doc",
		)
	}
	obj, ok := doc.Dict.Get(key)
	if !ok {
		lib.Die(
			node.Token,
			"key %s do not exist in %s",
			key.Inspect(),
			doc.Inspect(),
		)
	}

	return obj
}
