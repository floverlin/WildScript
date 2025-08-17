package enviroment

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func (e *Enviroment) loadBuiltin() {
	e.Set("print", &Func{
		Builtin: func(be blockEvaluator, args ...Object) Object {
			for _, arg := range args {
				if arg.Type() == OBJ_TYPE {
					obj := arg.(*Obj)
					f, ok := obj.Runes[STR_RUNE]
					if ok {
						runes := map[string]Object{
							SELF_RUNE: arg,
						}
						arg = be.EvalBlock(f.(*Func).Body, nil, runes)
					}
				}
				fmt.Print(arg.Inspect())
			}
			return Global[GLOBAL_NIL]
		},
	})

	e.Set("input", &Func{
		Builtin: func(_ blockEvaluator, args ...Object) Object {
			var value string
			fmt.Scanln(&value)
			return &Str{Value: value}
		},
	})

	e.Set("random", &Func{
		Builtin: func(_ blockEvaluator, args ...Object) Object {
			return &Num{Value: rand.Float64()}
		},
	})

	e.Set("type", &Func{
		Builtin: func(_ blockEvaluator, args ...Object) Object {
			obj := args[0]
			return &Str{Value: string(obj.Type())}
		},
	})

	e.Set("sleep", &Func{
		Builtin: func(_ blockEvaluator, args ...Object) Object {
			t := args[0].(*Num).Value
			ns := time.Duration(t * 1000 * 1000 * 1000)
			time.Sleep(ns * time.Nanosecond)
			return Global[GLOBAL_NIL]
		},
	})

	e.Set("len", &Func{Builtin: func(_ blockEvaluator, args ...Object) Object {
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
		case *List:
			return &Num{Value: float64(len(obj.Elements))}
		case *Obj:
			return &Num{Value: float64(len(obj.Fields))}
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
