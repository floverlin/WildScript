package enviroment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type ObjectType string

const (
	NIL  ObjectType = "nil"
	NUM  ObjectType = "num"
	STR  ObjectType = "str"
	BOOL ObjectType = "bool"
	DOC  ObjectType = "doc"
	FUNC ObjectType = "func"

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
	List  []Object
	Dict  *Dict
	Attrs map[string]Object
	Meta  *Doc
}

func NewDoc() *Doc {
	return &Doc{
		Attrs: make(map[string]Object),
		Dict:  NewDict(),
		Meta:  nil,
	}
}

func (d *Doc) Type() ObjectType { return DOC }
func (d *Doc) Inspect() string {
	var sb strings.Builder
	sb.WriteString("{")
	for key, elem := range d.Attrs {
		sb.WriteString(fmt.Sprintf("%s = %s, ", key, elem.Inspect()))
	}
	sb.WriteString(d.Dict.String())
	for idx, elem := range d.List {
		sb.WriteString(fmt.Sprintf("[%d]: %s, ", idx, elem.Inspect()))
	}
	result := sb.String()
	result = result[:len(result)-2]
	result += "}"
	return result
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
