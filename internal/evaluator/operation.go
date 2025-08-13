package evaluator

import (
	"errors"
	"math"
	"wildscript/internal/enviroment"
)

var numOps = map[string]func(left, right float64) (enviroment.Object, error){
	"+": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Num{Value: left + right}, nil
	},
	"-": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Num{Value: left - right}, nil
	},
	"*": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Num{Value: left * right}, nil
	},
	"/": func(left, right float64) (enviroment.Object, error) {
		if right == 0 {
			return nil, errors.New("division by zero")
		}
		return &enviroment.Num{Value: left / right}, nil
	},
	"//": func(left, right float64) (enviroment.Object, error) {
		if right == 0 {
			return nil, errors.New("division by zero")
		}
		return &enviroment.Num{Value: math.Floor(left / right)}, nil
	},
	"%": func(left, right float64) (enviroment.Object, error) {
		if right == 0 {
			return nil, errors.New("modulo by zero")
		}
		return &enviroment.Num{Value: math.Mod(left, right)}, nil
	},
	"^": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Num{Value: math.Pow(left, right)}, nil
	},

	"==": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left == right}, nil
	},
	"!=": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left != right}, nil
	},
	"<": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left < right}, nil
	},
	">": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left > right}, nil
	},
	"<=": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left <= right}, nil
	},
	">=": func(left, right float64) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left >= right}, nil
	},
}

var boolOps = map[string]func(left, right bool) (enviroment.Object, error){
	"==": func(left, right bool) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left == right}, nil
	},
	"!=": func(left, right bool) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left != right}, nil
	},
	"||": func(left, right bool) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left || right}, nil
	},
	"&&": func(left, right bool) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left && right}, nil
	},
}

var strOps = map[string]func(left, right string) (enviroment.Object, error){
	"+": func(left, right string) (enviroment.Object, error) {
		return &enviroment.Str{Value: left + right}, nil
	},

	"==": func(left, right string) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left == right}, nil
	},
	"!=": func(left, right string) (enviroment.Object, error) {
		return &enviroment.Bool{Value: left != right}, nil
	},
}
