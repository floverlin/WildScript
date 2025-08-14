package enviroment

import "fmt"

type methodMap = map[string]*Func

var typeMethodMap = map[ObjectType]methodMap{
	STR_TYPE: {
		"pr": &Func{
			Builtin: func(o ...Object) Object {
				fmt.Println("pr-pr")
				return &Nil{}
			},
		},
	},
}

// TODO ERROR RETURN
func FindMethod(objType ObjectType, name string) *Func {
	m, ok := typeMethodMap[objType]
	if !ok {
		return nil
	}
	f, ok := m[name]
	if !ok {
		return nil
	}
	return f
}
