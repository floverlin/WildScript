package enviroment

import (
	"fmt"
	"math"
)

func loadBuiltin(e *Enviroment) {
	e.Set("print", &Func{
		Builtin: func(args ...Object) Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())

			}
			fmt.Println()
			return &e.Single().Nil
		},
	})

	e.Set("type", &Func{
		Builtin: func(args ...Object) Object {
			obj := args[0]
			return &Str{Value: string(obj.Type())}
		},
	})

	e.Set("rune", &Func{
		Builtin: func(args ...Object) Object {
			obj := args[0]
			name := obj.(*Str).Value // MAY PANIC
			NewRune(name)
			return &e.Single().Nil
		},
	})

	e.Set("len", &Func{Builtin: func(args ...Object) Object {
		switch obj := args[0].(type) {
		case *Num:
			return &Num{Value: math.Floor(obj.Value)}

		case *Str:
			return &Num{Value: float64(len(obj.Value))}

		case *Bool:
			var result float64
			if obj.Value {
				result = 1
			}
			return &Num{Value: result}

		case *Nil:
			return &Num{Value: 0}

		case *Func:
			params := obj.LenOfParameters()
			return &Num{Value: float64(params)}
		default:
			panic(
				fmt.Sprintf(
					"unknown object type: %s",
					obj.Type(),
				),
			)
		}
	}})
}
