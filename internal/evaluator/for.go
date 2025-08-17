package evaluator

import (
	"wildscript/internal/ast"
	"wildscript/internal/enviroment"
)

func (e *Evaluator) evalForExpression(
	node *ast.ForExpression,
) enviroment.Object {
	var result enviroment.Object = enviroment.Global[enviroment.GLOBAL_NIL]

	cond := e.Eval(node.Condition)

	if cond.Type() == enviroment.BOOL_TYPE {
		idx := 0
		for {
			runes := Arguments{
				enviroment.IDX_RUNE: &enviroment.Num{Value: float64(idx)},
				enviroment.KEY_RUNE: enviroment.Global[enviroment.GLOBAL_NIL],
				enviroment.VAL_RUNE: enviroment.Global[enviroment.GLOBAL_TRUE],
			}
			result = e.EvalBlock(node.Body, nil, runes)
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
			var runes Arguments
			switch c := cond.(type) {
			case *enviroment.Num:
				runes = Arguments{
					"idx": &enviroment.Num{Value: float64(idx)},
					"key": &enviroment.Num{Value: float64(idx)},
					"val": &enviroment.Num{Value: float64(idx)},
				}
			case *enviroment.Str:
				runes = Arguments{
					"idx": &enviroment.Num{Value: float64(idx)},
					"key": &enviroment.Num{Value: float64(idx)},
					"val": &enviroment.Str{Value: string([]rune(c.Value)[idx])},
				}
			case *enviroment.Func:
				runes = Arguments{
					"idx": &enviroment.Num{Value: float64(idx)},
					"key": &enviroment.Str{Value: c.Parameters[idx].Value},
					"val": enviroment.Global[enviroment.GLOBAL_NIL],
				}
			case *enviroment.List:
				runes = Arguments{
					"idx": &enviroment.Num{Value: float64(idx)},
					"key": &enviroment.Num{Value: float64(idx)},
					"val": c.Elements[idx],
				}
			case *enviroment.Obj:
				vals := []enviroment.Object{}
				for _, val := range c.Fields {
					vals = append(vals, val)
				}
				runes = Arguments{
					"idx": &enviroment.Num{Value: float64(idx)},
					"key": enviroment.Global[enviroment.GLOBAL_NIL],
					"val": vals[idx],
				}
			default:
				panic("TODO")
			}

			result = e.EvalBlock(node.Body, nil, runes)

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
