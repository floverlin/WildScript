package enviroment

import (
	"fmt"
	"maps"
)

type method = func(self Object, be blockEvaluator, args ...Object) Object

type methodMap = map[string]method

var typeMethodMap = map[ObjectType]methodMap{
	STR_TYPE: {
		"wow": func(self Object, _ blockEvaluator, args ...Object) Object {
			selfStr := self.(*Str)

			fmt.Printf("WOW! %s :P\n", selfStr.Value)
			return &Nil{}
		},
	},

	LIST_TYPE: {
		"map": func(self Object, be blockEvaluator, args ...Object) Object {
			selfList := self.(*List)
			f := args[0].(*Func)

			newList := &List{
				Elements: []Object{},
			}

			paramName := f.Parameters[0].Value

			for _, elem := range selfList.Elements {
				newElem := be.EvalBlock(
					f.Body,
					map[string]Object{paramName: elem},
				)
				newList.Elements = append(newList.Elements, newElem)
			}

			return newList
		},
	},

	OBJ_TYPE: {
		"merge": func(self Object, _ blockEvaluator, args ...Object) Object {
			selfObj := self.(*Obj)
			otherObj := args[0].(*Obj)

			maps.Copy(selfObj.Fields, otherObj.Fields)
			maps.Copy(selfObj.Runes, otherObj.Runes)

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
		Builtin: func(be blockEvaluator, args ...Object) Object {
			return f(obj, be, args...)
		},
	}
}
