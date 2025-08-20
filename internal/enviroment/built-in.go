package enviroment

import (
	"fmt"
	"maps"
	"wildscript/internal/ast"
)

func (e *Enviroment) loadBuiltin() {
	e.Create("__print", &Func{
		Impl:   ast.FUNCTION,
		Native: print,
	})
	e.Create("print", &Func{
		Impl: ast.FUNCTION,
		Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			for idx, arg := range args {
				print(be, self, arg)
				if idx != len(args)-1 {
					fmt.Print(" ")
				}
			}
			return GLOBAL_NIL, nil
		},
	})
	e.Create("println", &Func{
		Impl: ast.FUNCTION,
		Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			for idx, arg := range args {
				print(be, self, arg)
				if idx != len(args)-1 {
					fmt.Print(" ")
				} else {
					fmt.Print("\n")
				}
			}
			return GLOBAL_NIL, nil
		},
	})
	e.Create("input", &Func{
		Impl: ast.FUNCTION,
		Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			var input string
			fmt.Scanln(&input)
			return &Str{Value: input}, nil
		},
	})
	e.Create("set_meta", &Func{
		Impl: ast.FUNCTION,
		Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			args[0].(*Doc).Meta = args[1].(*Doc)
			return GLOBAL_NIL, nil
		},
	})
	e.Create("get_meta", &Func{
		Impl: ast.FUNCTION,
		Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			return args[0].(*Doc).Meta, nil
		},
	})
	e.Create("merge", &Func{
		Impl: ast.FUNCTION,
		Native: func(be blockEvaluator, self Object, args ...Object) (Object, error) {
			left := args[0].(*Doc)
			right := args[1].(*Doc)
			maps.Copy(left.Attrs, right.Attrs)
			return GLOBAL_NIL, nil
		},
	})
}

func print(be blockEvaluator, self Object, args ...Object) (Object, error) {
	f, ok := lookupMeta(args[0], "__str")
	if !ok {
		fmt.Print(args[0].Type())
		return GLOBAL_NIL, nil
	}
	printable, err := f.(*Func).Call(be, args[0])
	if err != nil {
		return nil, err
	}
	fmt.Print(printable.Inspect())

	return GLOBAL_NIL, nil
}
