package enviroment

import (
	"strconv"

	"github.com/fatih/color"
)

type ObjectType string

const (
	NIL  ObjectType = "nil"
	NUM  ObjectType = "num"
	STR  ObjectType = "str"
	BOOL ObjectType = "bool"
	DOC  ObjectType = "doc"

	FUNCTION        ObjectType = "function"
	NATIVE_FUNCTION ObjectType = "native_function"

	SIGNAL ObjectType = "__signal"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Num struct {
	Value float64
}

func (f *Num) Type() ObjectType { return NUM }
func (f *Num) Inspect() string {
	return color.CyanString(
		strconv.FormatFloat(f.Value, 'g', -1, 64),
	)
}

type Str struct {
	Value string
}

func (s *Str) Type() ObjectType { return STR }
func (s *Str) Inspect() string {
	return s.Value
}

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType { return BOOL }
func (b *Bool) Inspect() string {
	return color.MagentaString(
		strconv.FormatBool(b.Value),
	)
}

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL }
func (n *Nil) Inspect() string {
	return color.BlueString("nil")
}

type Doc struct {
	List     []Object
	Dict     map[Object]Object
	Elements map[string]Object
}

func (d *Doc) Type() ObjectType { return DOC }
func (d *Doc) Inspect() string {
	return color.MagentaString("doc")
}

// ------------------------------  CONTROL  ------------------------------ //

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return SIGNAL }
func (r *Return) Inspect() string  { return "__return" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return SIGNAL }
func (c *Continue) Inspect() string  { return "__continue" }

type Break struct{}

func (b *Break) Type() ObjectType { return SIGNAL }
func (b *Break) Inspect() string  { return "__break" }
