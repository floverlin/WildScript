package enviroment

import (
	"errors"
	"fmt"
	"wildscript/internal/ast"

	"github.com/fatih/color"
)

type Native func(...Object) Object

type Func struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockExpression
	Enviroment *Enviroment
	Native     Native
	Impl       ast.FunctionImplementation
}

func (f *Func) Type() ObjectType { return FUNC }
func (f *Func) Inspect() string {
	return color.MagentaString(
		fmt.Sprintf("func<%s>", f.Impl))
}

type blockEvaluator interface {
	EvalBlock(
		block *ast.BlockExpression,
		outer *Enviroment,
		args map[string]Object,
	) Object
}

func (f *Func) Call(be blockEvaluator, args ...Object) (Object, error) {
	if len(args) != len(f.Parameters) {
		return nil, fmt.Errorf(
			"function want %d argument(s) got %d",
			len(f.Parameters),
			len(args),
		)
	}

	fArgs := map[string]Object{} // args
	for idx, arg := range args {
		fArgs[f.Parameters[idx].Value] = arg
	}

	result := be.EvalBlock(f.Body, f.Enviroment, fArgs)

	if result.Type() == SIGNAL {
		if ret, ok := result.(*Return); ok {
			return ret.Value, nil
		} else {
			return nil, errors.New("continue or break in function")
		}
	}

	return GLOBAL_NIL, nil
}
