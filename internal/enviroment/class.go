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
			Native: func(o ...Object) Object {
				self := o[0].(*Doc)
				self.List = append(self.List, o[1:]...)
				return self
			},
		},
		"reverse": &Func{
			Impl: ast.METHOD,
			Native: func(o ...Object) Object {
				self := o[0].(*Doc)
				slices.Reverse(self.List)
				return self
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
			Native: func(o ...Object) Object {
				self := o[0]
				fmt.Println("HOP!")
				return self
			},
		},
	},
}
