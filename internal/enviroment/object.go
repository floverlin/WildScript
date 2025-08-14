package enviroment

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	RUNE_TYPE ObjectType = "rune"
	LIST_TYPE ObjectType = "list"
	OBJ_TYPE  ObjectType = "obj"

	CONTROL_TYPE ObjectType = "CONTROL"
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
	return color.YellowString(s.Value)
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

type Evaluator interface {
	Eval(ast.Node, map[string]Object) Object
}

type Func struct {
	Builtin func(e Evaluator, args ...Object) Object

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

type List struct {
	Elements []Object
}

func (l *List) Type() ObjectType { return LIST_TYPE }
func (l *List) Inspect() string {
	var out strings.Builder
	out.WriteByte('[')
	for idx, elem := range l.Elements {
		out.WriteString(elem.Inspect())
		if idx != len(l.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteByte(']')
	return out.String()
}

type Obj struct {
	Fields map[string]Object
}

func (o *Obj) Type() ObjectType { return OBJ_TYPE }
func (o *Obj) Inspect() string {
	var out strings.Builder
	out.WriteString("{")
	idx := 0
	for key, value := range o.Fields {
		out.WriteString(
			fmt.Sprintf(
				"%s: %s",
				key,
				value.Inspect(),
			),
		)
		if idx != len(o.Fields)-1 {
			out.WriteString(", ")
		}
		idx++
	}
	out.WriteString("}")
	return out.String()
}

// ------------------------------  CONTROL  ------------------------------ //

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return CONTROL_TYPE }
func (r *Return) Inspect() string  { return "return" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTROL_TYPE }
func (c *Continue) Inspect() string  { return "continue" }
