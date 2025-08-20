package evaluator

import (
	"errors"
	"fmt"
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalAssignStatement(
	stmt *ast.AssignStatement,
) enviroment.Object {
	var err error
	var result enviroment.Object

	right := e.Eval(stmt.Right)
	switch left := stmt.Left.(type) {
	case *ast.Identifier:
		result, err = e.evalIdentifierAssign(left, right)
	case *ast.AttributeExpression:
		result, err = e.evalAttributeAssign(left, right)
	case *ast.IndexExpression:
		result, err = e.evalIndexAssign(left, right)
	case *ast.SliceExpression:
		result, err = e.evalSliceAssign(left, right)
	case *ast.KeyExpression:
		result, err = e.evalKeyAssign(left, right)
	}

	if err != nil {
		lib.Die(stmt.Token, err.Error())
	}

	return result
}

func (e *Evaluator) evalIdentifierAssign(
	left *ast.Identifier,
	value enviroment.Object,
) (enviroment.Object, error) {
	result, ok := e.env.Set(left.Value, value)
	if !ok {
		return nil, fmt.Errorf(
			"variable %s already exists",
			left.Value,
		)
	}
	return result, nil
}

func (e *Evaluator) evalAttributeAssign(
	left *ast.AttributeExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	object := e.Eval(left.Left)
	prop := &enviroment.Str{Value: left.Attribute.Value}

	f, err := lookupMeta(object, "__set_attribute")
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, object, prop, value)
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	return result, nil
}

func (e *Evaluator) evalIndexAssign(
	left *ast.IndexExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	object := e.Eval(left.Left)
	index := e.Eval(left.Index)

	if index.Type() != enviroment.NUM {
		return nil, errors.New("non num index type")
	}

	f, metaErr := lookupMeta(object, "__set_index")
	if metaErr != nil {
		lib.Die(
			left.Token,
			metaErr.Error(),
		)
	}
	result, err := f.Call(e, object, index, value)

	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	return result, nil
}

func (e *Evaluator) evalSliceAssign(
	left *ast.SliceExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	object := e.Eval(left.Left)
	start := e.Eval(left.Start)
	end := e.Eval(left.End)

	if (start.Type() != enviroment.NUM &&
		start.Type() != enviroment.NIL) ||
		(end.Type() != enviroment.NUM &&
			end.Type() != enviroment.NIL) {
		lib.Die(
			left.Token,
			"non num index",
		)
	}

	f, metaErr := lookupMeta(object, "__set_slice")
	if metaErr != nil {
		lib.Die(
			left.Token,
			metaErr.Error(),
		)
	}
	result, err := f.Call(e, object, start, end, value)

	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	return result, nil
}

func (e *Evaluator) evalKeyAssign(
	left *ast.KeyExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	object := e.Eval(left.Left)
	key := e.Eval(left.Key)

	if key.Type() == enviroment.NIL {
		f, err := lookupMeta(object, "__set_dict")
		if err != nil {
			lib.Die(
				left.Token,
				err.Error(),
			)
		}

		result, err := f.Call(e, object, value)
		if err != nil {
			lib.Die(
				left.Token,
				err.Error(),
			)
		}

		return result, nil
	}

	f, err := lookupMeta(object, "__set_key")
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	result, err := f.Call(e, object, key, value)
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	return result, nil
}

func (e *Evaluator) evalLetStatement(
	stmt *ast.LetStatement,
) enviroment.Object {
	right := e.Eval(stmt.Right)

	result, ok := e.env.Create(stmt.Left.Value, right)
	if !ok {
		lib.Die(
			stmt.Token,
			"variable %s already exists",
			stmt.Left.Value,
		)
	}
	return result
}
