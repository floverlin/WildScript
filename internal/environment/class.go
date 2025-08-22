package environment

import (
	"fmt"
	"slices"
)

func NewResult(value Object, ok *boolean) *document {
	r := NewDocument()
	r.Attrs["value"] = value
	r.Attrs["ok"] = ok
	return r
}

var iterMeta = func() *document {
	iter := NewDocument()
	iter.Attrs["__next"] = NewNative(func(
		be blockEvaluator,
		self Object,
		args ...Object,
	) (Object, error) {
		s := args[0].(*document)
		idx := int(s.Attrs["index"].(*number).Value)
		s.Attrs["index"] = NewNumber(float64(idx + 1))
		if idx >= len(s.List) {
			return NewResult(NewNil(), NewBoolean(false)), nil
		}
		return NewResult(s.List[idx], NewBoolean(true)), nil
	})
	return iter
}()

func UnpackResult(object Object) (Object, bool, error) {
	if doc, ok := object.(*document); ok {
		val, valOk := doc.Attrs["value"]
		ok, okOk := doc.Attrs["ok"]
		okBool, err := CheckBool(ok)
		if err != nil {
			return nil, false, err
		}
		if valOk && okOk {
			return val, okBool, nil
		}
	}
	return nil, false, fmt.Errorf("unpack result want document, gor %s", object.Type())
}

var classList = &document{
	List: []Object{},
	Dict: NewDict(),
	Attrs: map[string]Object{
		"append": NewNative(func(
			be blockEvaluator,
			self Object,
			args ...Object,
		) (Object, error) {
			s := args[0].(*document)
			s.List = append(s.List, args[1:]...)
			return s, nil
		}),
		"reverse": NewNative(func(
			be blockEvaluator,
			self Object,
			args ...Object,
		) (Object, error) {
			s := args[0].(*document)
			slices.Reverse(s.List)
			return s, nil
		}),
		"__iter": NewNative(func(
			be blockEvaluator,
			self Object,
			args ...Object,
		) (Object, error) {

			s := args[0].(*document)

			iter := NewDocument()
			iter.List = s.List
			iter.Attrs["index"] = NewNumber(0)
			iter.Meta = iterMeta
			return iter, nil
		}),
	},
}

var classDict = &document{
	List: []Object{},
	Dict: NewDict(),
	Attrs: map[string]Object{
		"hop": NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			s := args[0].(*document)
			fmt.Println("HOP!")
			return s, nil
		}),
	},
}
