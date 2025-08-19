package enviroment

import (
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
