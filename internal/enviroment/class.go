package enviroment

import (
	"fmt"
	"slices"
	"wildscript/internal/ast"
)

var classList = &Doc{
	List: []Object{},
	Dict: NewDict(),
	Attrs: map[string]Object{
		"append": &Func{
			Impl: ast.METHOD,
			Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
				s := self.(*Doc)
				s.List = append(s.List, args[1:]...)
				return s, nil
			},
		},
		"reverse": &Func{
			Impl: ast.METHOD,
			Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
				s := self.(*Doc)
				slices.Reverse(s.List)
				return s, nil
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
