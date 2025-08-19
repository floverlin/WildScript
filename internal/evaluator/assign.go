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
	case *ast.PropertyExpression:
		result, err = e.evalPropertyAssign(left, right)
	case *ast.IndexExpression:
		result, err = e.evalIndexAssign(left, right)
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

func (e *Evaluator) evalPropertyAssign(
	left *ast.PropertyExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	object := e.Eval(left.Left)
	prop := &enviroment.Str{Value: left.Property.Value}

	f, err := lookupMeta(object, "__set_property")
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	result, err := f(object, prop, value)
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

	f, err := lookupMeta(object, "__set_index")
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	result, err := f(object, index, value)
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

	f, err := lookupMeta(object, "__set_key")
	if err != nil {
		lib.Die(
			left.Token,
			err.Error(),
		)
	}

	result, err := f(object, key, value)
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
