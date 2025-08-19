package enviroment

import (
	"fmt"
	"wildscript/internal/ast"
)

func (e *Enviroment) loadBuiltin() {
	e.Create("__print", &Func{
		Impl:   ast.NATIVE,
		Native: print,
	})
	e.Create("print", &Func{
		Impl: ast.NATIVE,
		Native: func(o ...Object) Object {
			for idx, arg := range o {
				print(arg)
				if idx != len(o)-1 {
					fmt.Print(" ")
				}
			}
			return GLOBAL_NIL
		},
	})
	e.Create("println", &Func{
		Impl: ast.NATIVE,
		Native: func(o ...Object) Object {
			for idx, arg := range o {
				print(arg)
				if idx != len(o)-1 {
					fmt.Print(" ")
				} else {
					fmt.Print("\n")
				}
			}
			return GLOBAL_NIL
		},
	})
	e.Create("input", &Func{
		Impl: ast.NATIVE,
		Native: func(o ...Object) Object {
			var input string
			fmt.Scanln(&input)
			return &Str{Value: input}
		},
	})
}

func print(o ...Object) Object {
	fmt.Print(o[0].Inspect())
	return GLOBAL_NIL
}
