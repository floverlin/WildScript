package enviroment

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Callable interface {
	Call(be blockEvaluator, self Object, args ...Object) (Object, error)
}

type MetaFunc func(be blockEvaluator, self Object, args ...Object) (Object, error)

func (mf MetaFunc) Call(be blockEvaluator, self Object, args ...Object) (Object, error) {
	return mf(be, self, args...)
}

var DefaultMeta = map[ObjectType]map[string]MetaFunc{
	STR:  strMeta,
	BOOL: boolMeta,
	NUM:  numMeta,
	DOC:  docMeta,
	FUNC: funcMeta,
}

var boolMeta = map[string]MetaFunc{
	"__not": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if self.(*Bool).Value {
			return GLOBAL_FALSE, nil
		}
		return GLOBAL_TRUE, nil
	},
	"__eq": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Bool), args[0].(*Bool)
		if left.Value != right.Value {
			return GLOBAL_FALSE, nil
		}
		return GLOBAL_TRUE, nil
	},
	"__str": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if self.(*Bool).Value {
			return &Str{Value: "true"}, nil
		}
		return &Str{Value: "false"}, nil
	},
}

var funcMeta = map[string]MetaFunc{
	"__call": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		fmt.Println("!!")
		return self.(*Func).Call(be, self, args...)
	},
	"__str": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return &Str{Value: string(self.(*Func).Impl)}, nil
	},
}

var docMeta = map[string]MetaFunc{
	"__len": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		val := len(s.Attrs) + len(s.List) + s.Dict.Len()
		return &Num{Value: float64(val)}, nil
	},
	"__str": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return &Str{Value: "document"}, nil
	},
	"__bool": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		if s.Dict.Len() == 0 && len(s.List) == 0 && len(s.Attrs) == 0 {
			return GLOBAL_FALSE, nil
		}
		return GLOBAL_TRUE, nil
	},
	"__index": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		idx := int(args[0].(*Num).Value)
		if idx >= len(s.List) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		return s.List[idx], nil
	},
	"__set_index": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		idx := int(args[0].(*Num).Value)
		if idx >= len(s.List) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		s.List[idx] = args[0]
		return self, nil
	},
	"__key": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		result, ok := s.Dict.Get(args[0])
		if !ok {
			return nil, errors.New("key not exists")
		}
		return result, nil
	},
	"__set_key": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		s.Dict.Set(args[0], args[1])
		return self, nil
	},
	"__attribute": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		prop := args[0].(*Str)
		result, ok := s.Attrs[prop.Value]
		if !ok {
			return nil, errors.New("attribute not exists")
		}
		return result, nil
	},
	"__set_attribute": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		s := self.(*Doc)
		prop := args[0].(*Str)
		s.Attrs[prop.Value] = args[1]
		return self, nil
	},
}

var numMeta = map[string]MetaFunc{
	"__unm": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return &Num{Value: -self.(*Num).Value}, nil
	},
	"__add": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: left.Value + right.Value}, nil
	},
	"__sub": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: left.Value - right.Value}, nil
	},
	"__mul": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: left.Value * right.Value}, nil
	},
	"__div": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if right.Value == 0 {
			return nil, errors.New("division by zero")
		}
		return &Num{Value: left.Value / right.Value}, nil
	},
	"__floor_div": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if right.Value == 0 {
			return nil, errors.New("division by zero")
		}
		return &Num{Value: math.Floor(left.Value / right.Value)}, nil
	},
	"__mod": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if right.Value == 0 {
			return nil, errors.New("modulo by zero")
		}
		return &Num{Value: math.Mod(left.Value, right.Value)}, nil
	},
	"__pow": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		return &Num{Value: math.Pow(left.Value, right.Value)}, nil
	},
	"__str": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return &Str{
			Value: strconv.FormatFloat(
				self.(*Num).Value,
				'g', -1, 64,
			),
		}, nil
	},
	"__bool": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if self.(*Num).Value != 0 {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__eq": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value == right.Value {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__ne": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value != right.Value {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__lt": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value < right.Value {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__le": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value <= right.Value {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__gt": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value > right.Value {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
	"__ge": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Num), args[0].(*Num)
		if left.Value >= right.Value {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
}

var strMeta = map[string]MetaFunc{
	"__add": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Str), args[0].(*Str)
		return &Str{Value: left.Value + right.Value}, nil
	},
	"__eq": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Str), args[0].(*Str)
		return &Bool{Value: left.Value == right.Value}, nil
	},
	"__ne": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		left, right := self.(*Str), args[0].(*Str)
		return &Bool{Value: left.Value != right.Value}, nil
	},
	"__len": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		return &Num{Value: float64(len([]rune(self.(*Str).Value)))}, nil
	},
	"__index": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		sl := []rune(self.(*Str).Value)
		idx := int(args[0].(*Num).Value)
		if idx >= len(sl) || idx < 0 {
			return nil, errors.New("index out of range")
		}
		return &Str{Value: string(sl[idx])}, nil
	},
	"__slice": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		sl := []rune(self.(*Str).Value)
		var start, end int
		if _, ok := args[0].(*Nil); ok {
			start = 0
		} else {
			start = int(args[0].(*Num).Value)
		}
		if _, ok := args[1].(*Nil); ok {
			end = len(sl)
		} else {
			end = int(args[1].(*Num).Value)
		}
		if start < 0 || end > len(sl) {
			return nil, errors.New("index out of range")
		}
		return &Str{Value: string(sl[start:end])}, nil
	},
	"__num": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		result, err := strconv.ParseFloat(self.(*Str).Value, 64)
		if err != nil {
			return nil, err
		}
		return &Num{Value: result}, nil
	},
	"__bool": func(be blockEvaluator, self Object, args ...Object) (Object, error) {
		if len(self.(*Str).Value) != 0 {
			return GLOBAL_TRUE, nil
		}
		return GLOBAL_FALSE, nil
	},
}
