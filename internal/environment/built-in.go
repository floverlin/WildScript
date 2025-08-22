package environment

import (
	"fmt"
	"maps"
)

func (e *Environment) loadBuiltin() {
	e.Create("__print", NewNative(print))

	e.Create("print", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		for idx, arg := range args {
			print(be, self, arg)
			if idx != len(args)-1 {
				fmt.Print(" ")
			}
		}
		return NewNil(), nil
	}))

	e.Create("println", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		for idx, arg := range args {
			print(be, self, arg)
			if idx != len(args)-1 {
				fmt.Print(" ")
			} else {
				fmt.Print("\n")
			}
		}
		return NewNil(), nil
	}))

	e.Create("input", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		var input string
		fmt.Scanln(&input)
		return NewString(input), nil
	}))

	e.Create("set_meta", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		args[0].(*document).Meta = args[1].(*document)
		return NewNil(), nil
	}))

	e.Create("get_meta", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return args[0].(*document).Meta, nil
	}))

	e.Create("merge", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left := args[0].(*document)
		right := args[1].(*document)
		maps.Copy(left.Attrs, right.Attrs)
		return NewNil(), nil
	}))

	e.Create("str", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s, err := MetaCall(args[0], "__str", be, nil)
		if err != nil {
			return nil, fmt.Errorf("str: %w", err)
		}
		return s, nil
	}))

	e.Create("num", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		n, err := MetaCall(args[0], "__num", be, nil)
		if err != nil {
			return nil, fmt.Errorf("num: %s", err)
		}
		return n, nil
	}))

	e.Create("bool", NewNative(func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		b, err := MetaCall(args[0], "__bool", be, nil)
		if err != nil {
			return nil, fmt.Errorf("bool: %s", err)
		}
		return b, nil
	}))

}

func print(be blockEvaluator, self Object, args ...Object) (Object, error) {
	str, err := MetaCall(args[0], "__str", be, nil)
	if err != nil {
		return nil, fmt.Errorf("print: %w", err)
	}
	fmt.Print(str.Inspect())

	return NewNil(), nil
}
