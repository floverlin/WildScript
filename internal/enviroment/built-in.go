package enviroment

import (
	"fmt"
)

func (e *Enviroment) loadBuiltin() {
	e.Set("println", &NativeFunction{
		Native: func(be blockEvaluator, args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			fmt.Println()
			return GLOBAL_NIL
		},
	})
}
