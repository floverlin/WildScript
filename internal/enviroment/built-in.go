package enviroment

import (
	"fmt"
	"math"
	"reflect"
)

func loadBuiltin(e *Enviroment) {
	e.Set("print", &Func{
		Fn: func(args ...Object) Object {
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
		Fn: func(args ...Object) Object {
			obj := args[0]
			return &Str{Value: string(obj.Type())}
		},
	})

	e.Set("len", &Func{Fn: func(args ...Object) Object {
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
			args := reflect.ValueOf(obj.Fn).Type().NumIn()
			return &Num{Value: float64(args)}
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
