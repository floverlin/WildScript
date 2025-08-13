package enviroment

import (
	"fmt"
	"reflect"
	"strconv"
	"wildscript/internal/ast"

	"github.com/fatih/color"
)

type ObjectType string

const (
	NUM_TYPE  ObjectType = "num"
	STR_TYPE  ObjectType = "str"
	BOOL_TYPE ObjectType = "bool"
	NIL_TYPE  ObjectType = "nil"
	FUNC_TYPE ObjectType = "func"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Num struct {
	Value float64
}

func (f *Num) Type() ObjectType { return NUM_TYPE }
func (f *Num) Inspect() string {
	return color.GreenString(
		strconv.FormatFloat(f.Value, 'g', -1, 64),
	)
}

type Str struct {
	Value string
}

func (s *Str) Type() ObjectType { return STR_TYPE }
func (s *Str) Inspect() string {
	return color.YellowString(
		fmt.Sprintf("\"%s\"", s.Value),
	)
}

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType { return BOOL_TYPE }
func (b *Bool) Inspect() string {
	return color.MagentaString(
		strconv.FormatBool(b.Value),
	)
}

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_TYPE }
func (n *Nil) Inspect() string {
	return color.BlueString("nil")
}

// TODO
type Func struct {
	Builtin func(...Object) Object

	Parameters []*ast.Identifier
	Body       *ast.BlockExpression
	Enviroment *Enviroment
}

func (f *Func) Type() ObjectType { return FUNC_TYPE }
func (f *Func) Inspect() string {
	return color.CyanString("func")
}
func (f *Func) LenOfParameters() int {
	if f.Builtin != nil {
		return reflect.ValueOf(f.Builtin).Type().NumIn()
	}
	return len(f.Parameters)
}
