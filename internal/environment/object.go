package environment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type ObjectType string

const (
	NIL      ObjectType = "nil"
	NUMBER   ObjectType = "number"
	STRING   ObjectType = "string"
	BOOLEAN  ObjectType = "boolean"
	DOCUMENT ObjectType = "document"
	FUNCTION ObjectType = "function"

	SIGNAL ObjectType = "__signal"
)

var (
	globalFalse = &boolean{Value: false}
	globalTrue  = &boolean{Value: true}
	globalNil   = &nil_{}
)

func CheckBool(b Object) (bool, error) {
	if b, ok := b.(*boolean); ok {
		return b == globalTrue, nil
	}
	return false, fmt.Errorf("not bool value %s", b.Type())
}

func lookupDocMeta(doc *document, metaName string) (Object) {
	if result, ok := doc.Attrs[metaName]; ok {
		return result
	}
	if doc.Meta != nil {
		if result := lookupDocMeta(doc.Meta, metaName); result != nil {
			return result
		}
	}
	return nil
}

func MetaCall(
	object Object,
	metaName string,
	be blockEvaluator,
	self Object,
	args ...Object,
) (Object, error) {
	var metaFunc Object
	if doc, ok := object.(*document); ok {
		if doc.Meta != nil {
			metaFunc = lookupDocMeta(doc.Meta, metaName)
		}
	}

	if metaFunc == nil {
		typeMetaFuncs, ok := defaultMeta[object.Type()]
		if !ok {
			return nil, fmt.Errorf("type %s default meta not supported", object.Type())
		}
		result, ok := typeMetaFuncs[metaName]
		if !ok {
			return nil, fmt.Errorf("no %s method in %s", metaName, object.Type())
		}
		metaFunc = result
	}

	var result Object
	var err error
	if self != nil {
		args = append([]Object{self}, args...)
	}
	if mf, ok := metaFunc.(*function); ok {
		result, err = mf.Call(be, object, args...)
	} else {
		return nil, fmt.Errorf("%s not a function", metaFunc.Type())
	}

	if err != nil {
		return nil, fmt.Errorf("%s call: %w", metaFunc.Type(), err)
	}

	return result, err
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type number struct {
	Value float64
}

func NewNumber(value float64) *number {
	return &number{
		Value: value,
	}
}

func (f *number) Type() ObjectType { return NUMBER }
func (f *number) Inspect() string {
	return color.GreenString(
		strconv.FormatFloat(f.Value, 'g', -1, 64),
	)
}

type string_ struct {
	Value string
}

func NewString(value string) *string_ {
	return &string_{Value: value}
}

func (s *string_) Type() ObjectType { return STRING }
func (s *string_) Inspect() string {
	return s.Value
}

type boolean struct {
	Value bool
}

func NewBoolean(value bool) *boolean {
	if value {
		return globalTrue
	}
	return globalFalse
}

func (b *boolean) Type() ObjectType { return BOOLEAN }
func (b *boolean) Inspect() string {
	return color.MagentaString(
		strconv.FormatBool(b.Value),
	)
}

type nil_ struct{}

func NewNil() *nil_ {
	return globalNil
}

func (n *nil_) Type() ObjectType { return NIL }
func (n *nil_) Inspect() string {
	return color.BlueString("nil")
}

type document struct {
	List  []Object
	Dict  *Dict
	Attrs map[string]Object
	Meta  *document
}

func NewDocument() *document {
	return &document{
		Attrs: make(map[string]Object),
		Dict:  NewDict(),
		Meta:  nil,
	}
}

func (d *document) Type() ObjectType { return DOCUMENT }
func (d *document) Inspect() string {
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
	if len(result) > 1 {
		result = result[:len(result)-2]
	}
	result += "}"
	return result
}

// ------------------------------  CONTROL  ------------------------------ //

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return SIGNAL }
func (r *Return) Inspect() string  { return "__return" }

type Export struct {
	Value Object
}

func (e *Export) Type() ObjectType { return SIGNAL }
func (e *Export) Inspect() string  { return "__export" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return SIGNAL }
func (c *Continue) Inspect() string  { return "__continue" }

type Break struct{}

func (b *Break) Type() ObjectType { return SIGNAL }
func (b *Break) Inspect() string  { return "__break" }
