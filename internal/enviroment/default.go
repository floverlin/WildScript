package enviroment

import (
	"errors"
	"math"
	"strconv"
)

type MetaFunc func(self Object, args ...Object) (Object, error)

var DefaultMeta = map[ObjectType]map[string]MetaFunc{
	STR:  strMeta,
	BOOL: boolMeta,
	NUM:  numMeta,
	DOC:  docMeta,
}

var boolMeta = map[string]MetaFunc{
	"__eq": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Bool), args[0].(*Bool)
		if left.Value != right.Value {
			return GLOBAL_FALSE, nil
		}
		return GLOBAL_TRUE, nil
	},
	"__str": func(self Object, args ...Object) (Object, error) {
		if self.(*Bool).Value {
			return &Str{Value: "false"}, nil
		}
		return &Str{Value: "true"}, nil
	},
}

var docMeta = map[string]MetaFunc{
	"__len": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		val := len(s.Elements) + len(s.List) + s.Dict.Len()
		return &Num{Value: float64(val)}, nil
	},
	"__str": func(self Object, args ...Object) (Object, error) {
		return &Str{Value: "document"}, nil
	},
	"__index": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		idx := int(args[0].(*Num).Value)
		if idx >= len(s.List) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		return s.List[idx], nil
	},
	"__set_index": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		idx := int(args[0].(*Num).Value)
		if idx >= len(s.List) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		s.List[idx] = args[0]
		return args[0], nil
	},
	"__key": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		result, ok := s.Dict.Get(args[0])
		if !ok {
			return nil, errors.New("key not exists")
		}
		return result, nil
	},
	"__set_key": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		s.Dict.Set(args[0], args[1])
		return args[1], nil
	},
	"__property": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		prop := args[0].(*Str)
		result, ok := s.Elements[prop.Value]
		if !ok {
			return nil, errors.New("property not exists")
		}
		return result, nil
	},
	"__set_property": func(self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		prop := args[0].(*Str)
		s.Elements[prop.Value] = args[1]
		return args[1], nil
	},
}

var numMeta = map[string]MetaFunc{
	"__add": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: left.Value + right.Value}, nil
	},
	"__sub": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: left.Value - right.Value}, nil
	},
	"__mul": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: left.Value * right.Value}, nil
	},
	"__div": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if right.Value == 0 {
			return nil, errors.New("division by zero")
		}
		return &Num{Value: left.Value / right.Value}, nil
	},
	"__floor_div": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if right.Value == 0 {
			return nil, errors.New("division by zero")
		}
		return &Num{Value: math.Floor(left.Value / right.Value)}, nil
	},
	"__mod": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if right.Value == 0 {
			return nil, errors.New("modulo by zero")
		}
		return &Num{Value: math.Mod(left.Value, right.Value)}, nil
	},
	"__pow": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: math.Pow(left.Value, right.Value)}, nil
	},
	"__str": func(self Object, args ...Object) (Object, error) {
		return &Str{
			Value: strconv.FormatFloat(
				self.(*Num).Value,
				'g', -1, 64,
			),
		}, nil
	},
	"__bool": func(self Object, args ...Object) (Object, error) {
		if self.(*Num).Value != 0 {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__eq": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value != right.Value {
			return GLOBAL_FALSE, nil
		}
		return GLOBAL_TRUE, nil
	},
	"__lt": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value >= right.Value {
			return GLOBAL_FALSE, nil
		}
		return GLOBAL_TRUE, nil
	},
}

var strMeta = map[string]MetaFunc{
	"__add": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Str), args[0].(*Str)
		return &Str{Value: left.Value + right.Value}, nil
	},
	"__eq": func(self Object, args ...Object) (Object, error) {
		left, right := self.(*Str), args[0].(*Str)
		return &Bool{Value: left.Value == right.Value}, nil
	},
	"__len": func(self Object, args ...Object) (Object, error) {
		return &Num{Value: float64(len([]rune(self.(*Str).Value)))}, nil
	},
	"__index": func(self Object, args ...Object) (Object, error) {
		sl := []rune(self.(*Str).Value)
		idx := int(args[0].(*Num).Value)
		if idx >= len(sl) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		return &Str{Value: string(sl[idx])}, nil
	},
	"__slice": func(self Object, args ...Object) (Object, error) {
		sl := []rune(self.(*Str).Value)
		start := int(args[0].(*Num).Value)
		end := int(args[1].(*Num).Value)
		if start < 0 || end > len(sl) {
			return nil, errors.New("index out of range")
		}
		return &Str{Value: string(sl[start:end])}, nil
	},
	"__num": func(self Object, args ...Object) (Object, error) {
		result, err := strconv.ParseFloat(self.(*Str).Value, 64)
		if err != nil {
			return nil, err
		}
		return &Num{Value: result}, nil
	},
	"__bool": func(self Object, args ...Object) (Object, error) {
		if len(self.(*Str).Value) != 0 {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
}
