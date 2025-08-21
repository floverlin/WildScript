package enviroment

import (
	"fmt"
	"slices"
	"wildscript/internal/ast"
)

func NewResult(value Object, ok *Bool) *Doc {
	r := NewDoc()
	r.Attrs["value"] = value
	r.Attrs["ok"] = ok
	return r
}

var classList = &Doc{
	List: []Object{},
	Dict: NewDict(),
	Attrs: map[string]Object{
		"append": &Func{
			Impl: ast.METHOD,
			Native: func(
				be blockEvaluator,
				self Object,
				args ...Object,
			) (Object, error) {
				s := self.(*Doc)
				s.List = append(s.List, args[1:]...)
				return s, nil
			},
		},
		"reverse": &Func{
			Impl: ast.METHOD,
			Native: func(
				be blockEvaluator,
				self Object,
				args ...Object,
			) (Object, error) {
				s := self.(*Doc)
				slices.Reverse(s.List)
				return s, nil
			},
		},
		"__iter": &Func{
			Impl: ast.METHOD,
			Native: func(
				be blockEvaluator,
				self Object,
				args ...Object,
			) (Object, error) {
				s := self.(*Doc)

				iter := NewDoc()
				iter.List = s.List
				iter.Attrs["index"] = &Num{Value: 0}
				iter.Attrs["__next"] = &Func{
					Impl: ast.METHOD,
					Native: func(
						be blockEvaluator,
						self Object,
						args ...Object,
					) (Object, error) {
						s := self.(*Doc)
						idx := int(s.Attrs["index"].(*Num).Value)
						s.Attrs["index"] = &Num{Value: float64(idx + 1)}
						if idx >= len(s.List) {
							return NewResult(GLOBAL_NIL, GLOBAL_FALSE), nil
						}
						return NewResult(s.List[idx], GLOBAL_TRUE), nil
					},
				}
				return iter, nil
			},
		},
	},
}

var classDict = &Doc{
	List: []Object{},
	Dict: NewDict(),
	Attrs: map[string]Object{
		"hop": &Func{
			Impl: ast.METHOD,
			Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
				s := self.(*Doc)
				fmt.Println("HOP!")
				return s, nil
			},
		},
	},
}
