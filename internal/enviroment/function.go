package enviroment

import (
	"wildscript/internal/ast"

	"github.com/fatih/color"
)

type blockEvaluator interface {
	EvalBlock(
		*ast.BlockExpression,
		*Enviroment,
		map[string]Object,
	) Object
}

type Callable interface {
	Call(be blockEvaluator, args ...Object) Object
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockExpression
	Enviroment *Enviroment
	Impl       ast.FunctionImplementation
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	return color.MagentaString("function")
}

type NativeFunction struct {
	Native func(blockEvaluator, ...Object) Object
}

func (nf *NativeFunction) Type() ObjectType { return NATIVE_FUNCTION }
func (nf *NativeFunction) Inspect() string {
	return color.MagentaString("native_function")
}
