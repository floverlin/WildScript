package enviroment

import "fmt"

type method = func(self Object, args ...Object) Object

type methodMap = map[string]method

var typeMethodMap = map[ObjectType]methodMap{
	STR_TYPE: {
		"wow": func(self Object, args ...Object) Object {
			selfStr := self.(*Str)

			fmt.Printf("WOW! %s :P\n", selfStr.Value)
			return &Nil{}
		},
	},
}

// TODO ERROR RETURN
func FindMethod(obj Object, name string) *Func {
	m, ok := typeMethodMap[obj.Type()]
	if !ok {
		return nil
	}
	f, ok := m[name]
	if !ok {
		return nil
	}

	return &Func{
		Builtin: func(args ...Object) Object {
			return f(obj, args...)
		},
	}
}
