package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
)

func setForRunes(idx int, key, val enviroment.Object) {
	enviroment.TakeRune(enviroment.IDX_RUNE).
		Set(&enviroment.Num{Value: float64(idx)})

	enviroment.TakeRune(enviroment.KEY_RUNE).
		Set(key)

	enviroment.TakeRune(enviroment.VAL_RUNE).
		Set(val)
}

func (e *Evaluator) evalForExpression(
	node *ast.ForExpression,
) enviroment.Object {
	var result enviroment.Object = &e.env.Single().Nil

	cond := e.Eval(node.Condition)

	if cond.Type() == enviroment.BOOL_TYPE {
		idx := 0
		for {
			setForRunes(
				idx,
				&enviroment.Nil{},
				&enviroment.Bool{Value: true},
			)
			result = e.Eval(node.Body)
			cond = e.Eval(node.Condition)
			if !cond.(*enviroment.Bool).Value {
				break
			}
			idx++
		}

	} else {
		var iters int
		switch c := cond.(type) {
		case *enviroment.Num:
			iters = int(c.Value)
		case *enviroment.Str:
			iters = len(c.Value)
		case *enviroment.Nil:
			iters = 0
		case *enviroment.Func:
			iters = c.LenOfParameters()
		case *enviroment.List:
			iters = len(c.Elements)
		case *enviroment.Obj:
			iters = len(c.Fields)
		default:
			panic("TODO")
		}

		for idx := range iters {
			switch c := cond.(type) {
			case *enviroment.Num:
				setForRunes(
					idx,
					&enviroment.Num{Value: float64(idx)},
					&enviroment.Num{Value: float64(idx)},
				)
			case *enviroment.Str:
				setForRunes(
					idx,
					&enviroment.Num{Value: float64(idx)},
					&enviroment.Str{Value: string([]rune(c.Value)[idx])},
				)
			case *enviroment.Func:
				setForRunes(
					idx,
					&enviroment.Str{Value: c.Parameters[idx].Value},
					&enviroment.Nil{},
				)
			case *enviroment.List:
				setForRunes(
					idx,
					&enviroment.Num{Value: float64(idx)},
					c.Elements[idx],
				)
			case *enviroment.Obj:
				vals := []enviroment.Object{}
				for _, val := range c.Fields {
					vals = append(vals, val)
				}
				setForRunes(
					idx,
					&enviroment.Nil{},
					vals[idx],
				)
			default:
				panic("TODO")
			}

			result = e.EvalBlock(node.Body, nil)

			if result.Type() == enviroment.CONTROL_TYPE {
				switch res := result.(type) {
				case *enviroment.Continue:
					continue
				case *enviroment.Return:
					result = res.Value
				}
				break
			}
		}
	}

	return result
}
