package enviroment

import (
	"fmt"
	"math"
)

func loadBuiltin(e *Enviroment) {
	e.Set("print", &Func{
		Builtin: func(ev Evaluator, args ...Object) Object {
			for i, arg := range args {
				if arg.Type() == OBJ_TYPE {
					obj := arg.(*Obj)
					r := NewRune("str")
					f, ok := obj.Runes[r.ID]
					if ok {
						NewRune("self").Set(arg)
						arg = ev.Eval(f.(*Func).Body, nil)
					}
				}
				fmt.Print(arg.Inspect())
				if i != len(args)-1 {
					fmt.Print(" ")
				}
			}
			fmt.Println()
			return &e.Single().Nil
		},
	})

	e.Set("type", &Func{
		Builtin: func(_ Evaluator, args ...Object) Object {
			obj := args[0]
			return &Str{Value: string(obj.Type())}
		},
	})

	e.Set("rune", &Func{
		Builtin: func(_ Evaluator, args ...Object) Object {
			obj := args[0]
			name := obj.(*Str).Value // MAY PANIC
			NewRune(name)
			return &e.Single().Nil
		},
	})

	e.Set("len", &Func{Builtin: func(_ Evaluator, args ...Object) Object {
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
