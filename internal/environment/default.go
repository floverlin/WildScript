package environment

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"wildscript/internal/ast"
)

func LookupAttr(doc *document, attr string) (Object, bool) {
	if result, ok := doc.Attrs[attr]; ok {
		return result, ok
	}
	if doc.Meta != nil {
		if result, ok := LookupAttr(doc.Meta, attr); ok {
			return result, ok
		}
	}
	return nil, false
}

type Callable interface {
	Call(be blockEvaluator, self Object, args ...Object) (Object, error)
}

var defaultMeta = map[ObjectType]map[string]*function{
	STRING:   strMeta,
	BOOLEAN:  boolMeta,
	NUMBER:   numMeta,
	DOCUMENT: docMeta,
	FUNCTION: funcMeta,
}

var boolMeta = map[string]*function{
	"__not": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if self.(*boolean).Value {
			return NewBoolean(false), nil
		}
		return NewBoolean(true), nil
	}),
	"__eq": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*boolean), args[0].(*boolean)
		if left.Value != right.Value {
			return NewBoolean(false), nil
		}
		return NewBoolean(true), nil
	}),
	"__str": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if self.(*boolean).Value {
			return NewString("true"), nil
		}
		return NewString("false"), nil
	}),
}

var funcMeta = map[string]*function{
	"__call": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*function)
		if s.Impl == ast.METHOD {
			self = args[0]
			args = args[1:]
		}
		return s.Call(be, self, args...)
	}),
	"__str": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return NewString(string(self.(*function).Impl)), nil
	}),
}

var docMeta = map[string]*function{
	"__len": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		val := len(s.Attrs) + len(s.List) + s.Dict.Len()
		return NewNumber(float64(val)), nil
	}),
	"__str": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return NewString("document"), nil
	}),
	"__bool": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		if s.Dict.Len() == 0 && len(s.List) == 0 && len(s.Attrs) == 0 {
			return NewBoolean(false), nil
		}
		return NewBoolean(true), nil
	}),
	"__index": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		idx := int(args[0].(*number).Value)
		if idx >= len(s.List) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		return s.List[idx], nil
	}),
	"__set_index": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		idx := int(args[0].(*number).Value)
		if idx >= len(s.List) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		s.List[idx] = args[1]
		return self, nil
	}),
	"__key": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		result, ok := s.Dict.Get(args[0])
		if !ok {
			return nil, errors.New("key not exists")
		}
		return result, nil
	}),
	"__set_key": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		s.Dict.Set(args[0], args[1])
		return self, nil
	}),
	"__dict": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		dict := NewDocument()
		dict.Dict = s.Dict.Clone()
		dict.Meta = classDict
		return dict, nil
	}),
	"__set_dict": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		dict := args[0].(*document)
		s.Dict = dict.Dict.Clone()
		return self, nil
	}),
	"__attribute": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		prop := args[0].(*string_)
		if result, ok := LookupAttr(s, prop.Value); ok {
			return result, nil
		}
		return nil, errors.New("attribute not exists")
	}),
	"__set_attribute": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		prop := args[0].(*string_)
		s.Attrs[prop.Value] = args[1]
		return self, nil
	}),
	"__slice": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		var start, end int
		if _, ok := args[0].(*nil_); ok {
			start = 0
		} else {
			start = int(args[0].(*number).Value)
		}
		if _, ok := args[1].(*nil_); ok {
			end = len(s.List)
		} else {
			end = int(args[1].(*number).Value)
		}
		slice := NewDocument()
		slice.Meta = classList
		slice.List = slices.Clone(s.List[start:end])
		return slice, nil
	}),
	"__set_slice": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*document)
		var start, end int
		if _, ok := args[0].(*nil_); ok {
			start = 0
		} else {
			start = int(args[0].(*number).Value)
		}
		if _, ok := args[1].(*nil_); ok {
			end = len(s.List)
		} else {
			end = int(args[1].(*number).Value)
		}
		list := args[2].(*document).List
		list = append(list, s.List[end:]...)
		list = append(s.List[:start], list...)
		s.List = slices.Clone(list)
		return self, nil
	}),
}

var numMeta = map[string]*function{
	"__unm": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return NewNumber(-self.(*number).Value), nil
	}),
	"__add": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		return NewNumber(left.Value + right.Value), nil
	}),
	"__sub": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		return NewNumber(left.Value - right.Value), nil
	}),
	"__mul": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		return NewNumber(left.Value * right.Value), nil
	}),
	"__div": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if right.Value == 0 {
			return nil, errors.New("division by zero")
		}
		return NewNumber(left.Value / right.Value), nil
	}),
	"__floor_div": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if right.Value == 0 {
			return nil, errors.New("division by zero")
		}
		return NewNumber(math.Floor(left.Value / right.Value)), nil
	}),
	"__mod": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if right.Value == 0 {
			return nil, errors.New("modulo by zero")
		}
		return NewNumber(math.Mod(left.Value, right.Value)), nil
	}),
	"__pow": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		return NewNumber(math.Pow(left.Value, right.Value)), nil
	}),
	"__str": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return NewString(
			strconv.FormatFloat(
				self.(*number).Value,
				'g', -1, 64,
			),
		), nil
	}),
	"__bool": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if self.(*number).Value != 0 {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
	"__eq": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if left.Value == right.Value {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
	"__ne": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if left.Value != right.Value {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
	"__lt": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if left.Value < right.Value {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
	"__le": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if left.Value <= right.Value {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
	"__gt": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if left.Value > right.Value {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
	"__ge": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*number), args[0].(*number)
		if left.Value >= right.Value {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
}

var strMeta = map[string]*function{
	"__str": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return self, nil
	}),
	"__add": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*string_), args[0].(*string_)
		return NewString(left.Value + right.Value), nil
	}),
	"__eq": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*string_), args[0].(*string_)
		return NewBoolean(left.Value == right.Value), nil
	}),
	"__ne": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*string_), args[0].(*string_)
		return NewBoolean(left.Value != right.Value), nil
	}),
	"__len": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return NewNumber(float64(len([]rune(self.(*string_).Value)))), nil
	}),
	"__index": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		sl := []rune(self.(*string_).Value)
		idx := int(args[0].(*number).Value)
		if idx >= len(sl) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		return NewString(string(sl[idx])), nil
	}),
	"__slice": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		sl := []rune(self.(*string_).Value)
		var start, end int
		if _, ok := args[0].(*nil_); ok {
			start = 0
		} else {
			start = int(args[0].(*number).Value)
		}
		if _, ok := args[1].(*nil_); ok {
			end = len(sl)
		} else {
			end = int(args[1].(*number).Value)
		}
		if start < 0 || end > len(sl) {
			return nil, errors.New("index out of range")
		}
		return NewString(string(sl[start:end])), nil
	}),
	"__num": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		result, err := strconv.ParseFloat(self.(*string_).Value, 64)
		if err != nil {
			return nil, err
		}
		return NewNumber(result), nil
	}),
	"__bool": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if len(self.(*string_).Value) != 0 {
			return NewBoolean(true), nil
		}
		return NewBoolean(false), nil
	}),
}
