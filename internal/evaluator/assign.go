package evaluator

import (
	"errors"
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
	"wildscript/internal/logger"
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
	case *ast.PropertyAccessExpression:
		result, err = e.evalPropertyAssign(left, right)
	case *ast.IndexExpression:
		result, err = e.evalIndexAssign(left, right)
	}

	if err != nil {
		panic(
			logger.Slog(
				stmt.Token.Line,
				stmt.Token.Column,
				"%s", err.Error(),
			),
		)
	}

	return result
}

func (e *Evaluator) evalIdentifierAssign(
	left *ast.Identifier,
	value enviroment.Object,
) (enviroment.Object, error) {
	if left.IsRune {
		r, ok := enviroment.FindRune(left.Value)
		if !ok {
			return nil, errors.New("undefined rune")
		}
		result := r.Set(value)
		return result, nil
	} else if left.IsOuter {
		result, ok := e.env.SetOuter(left.Value, value)
		if !ok {
			return nil, errors.New("undefined variable")
		}
		return result, nil
	} else {
		result := e.env.Set(left.Value, value)
		return result, nil
	}
}

func (e *Evaluator) evalPropertyAssign(
	left *ast.PropertyAccessExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	obj := e.Eval(left.Object)
	prop := left.Property.Value
	if obj.Type() != enviroment.OBJ_TYPE {
		return nil, errors.New("assign property to non obj type")
	}
	obj.(*enviroment.Obj).Fields[prop] = value
	return obj, nil
}

func (e *Evaluator) evalIndexAssign(
	left *ast.IndexExpression,
	value enviroment.Object,
) (enviroment.Object, error) {
	list := e.Eval(left.Left)
	indexObj := e.Eval(left.Index)
	if indexObj.Type() != enviroment.NUM_TYPE {
		return nil, errors.New("non num index type")
	}
	idx := int(indexObj.(*enviroment.Num).Value)
	if list.Type() != enviroment.LIST_TYPE {
		return nil, errors.New("assign index to non list type")
	}
	list.(*enviroment.List).Elements[idx] = value
	return list, nil
}
