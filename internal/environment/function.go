package environment

import (
	"fmt"
	"wildscript/internal/ast"

	"github.com/fatih/color"
)

type Native func(
	be blockEvaluator,
	self Object,
	args ...Object,
) (Object, error)

type function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockExpression
	Environment *Environment
	Native      Native
	Impl        ast.FunctionImplementation
}

func NewNative(f Native) *function {
	return &function{
		Native: f,
	}
}

func NewFunction(params []*ast.Identifier,
	body *ast.BlockExpression,
	env *Environment,
	Impl ast.FunctionImplementation,
) *function {
	return &function{
		Parameters:  params,
		Body:        body,
		Environment: env,
		Impl:        Impl,
	}
}

func (f *function) Type() ObjectType { return FUNCTION }
func (f *function) Inspect() string {
	if f.Native != nil {
		return color.MagentaString("function<native>")
	}
	return color.MagentaString(
		fmt.Sprintf("function<%s>", f.Impl))
}

type blockEvaluator interface {
	EvalBlock(
		block *ast.BlockExpression,
		outer *Environment,
		args map[string]Object,
	) Object
}

func (f *function) Call(
	be blockEvaluator,
	self Object,
	args ...Object,
) (Object, error) {
	if f.Native != nil {
		return f.Native(be, self, args...)
	}

	if f.Impl == ast.METHOD {
		args = append([]Object{self}, args...)
	}

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

	result := be.EvalBlock(f.Body, f.Environment, fArgs)

	if result.Type() == SIGNAL {
		if ret, ok := result.(*Return); ok {
			return ret.Value, nil
		} else {
			return nil, fmt.Errorf("%s in function", result.Inspect())
		}
	}

	return NewNil(), nil
}
