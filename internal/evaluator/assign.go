package evaluator

import (
	"errors"
	"fmt"
	"wildscript/internal/ast"
	"wildscript/internal/environment"
	"wildscript/internal/lib"
)

func (e *Evaluator) evalAssignStatement(
	stmt *ast.AssignStatement,
) environment.Object {
	var err error
	var result environment.Object

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
	value environment.Object,
) (environment.Object, error) {
	result, ok := e.env.Set(left.Value, value)
	if !ok {
		return nil, fmt.Errorf(
			"variable %s not exists",
			left.Value,
		)
	}
	return result, nil
}

func (e *Evaluator) evalAttributeAssign(
	left *ast.AttributeExpression,
	value environment.Object,
) (environment.Object, error) {
	object := e.Eval(left.Left)
	prop := environment.NewString(left.Attribute.Value)

	result, err := environment.MetaCall(object, "__set_attribute", e, nil, prop, value)
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
	value environment.Object,
) (environment.Object, error) {
	object := e.Eval(left.Left)
	index := e.Eval(left.Index)

	if index.Type() != environment.NUMBER &&
		index.Type() != environment.NIL {
		return nil, errors.New("non num index type")
	}

	if index.Type() == environment.NIL {
		result, err := environment.MetaCall(object, "__set_list", e, nil, value)
		if err != nil {
			lib.Die(
				left.Token,
				err.Error(),
			)
		}

		return result, nil
	}

	result, err := environment.MetaCall(object, "__set_index", e, nil, index, value)
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
	value environment.Object,
) (environment.Object, error) {
	object := e.Eval(left.Left)
	start := e.Eval(left.Start)
	end := e.Eval(left.End)

	if (start.Type() != environment.NUMBER &&
		start.Type() != environment.NIL) ||
		(end.Type() != environment.NUMBER &&
			end.Type() != environment.NIL) {
		lib.Die(
			left.Token,
			"non num index",
		)
	}

	result, err := environment.MetaCall(object, "__set_slice", e, nil, start, end, value)
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
	value environment.Object,
) (environment.Object, error) {
	object := e.Eval(left.Left)
	key := e.Eval(left.Key)

	if key.Type() == environment.NIL {
		result, err := environment.MetaCall(object, "__set_dict", e, nil, value)
		if err != nil {
			lib.Die(
				left.Token,
				err.Error(),
			)
		}

		return result, nil
	}

	result, err := environment.MetaCall(object, "__set_key", e, nil, key, value)
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
) environment.Object {
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
