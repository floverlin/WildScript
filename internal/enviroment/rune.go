package enviroment

import (
	"fmt"

	"github.com/fatih/color"
)

var runeCounter uint64 = 0
var runeMap = map[uint64]Object{}

type Rune struct {
	id uint64
}

func (r *Rune) Type() ObjectType { return RUNE_TYPE }
func (r *Rune) Inspect() string {
	return color.CyanString(
		fmt.Sprintf(
			"rune<%d>",
			r.id,
		),
	)
}

func NewRune() *Rune {
	r := &Rune{
		id: runeCounter,
	}
	runeCounter++
	runeMap[r.id] = &Nil{}
	return r
}

func (r *Rune) Get() Object {
	obj, ok := runeMap[r.id]
	if !ok {
		panic("rune value is empty")
	}
	return obj
}

func (r *Rune) Set(obj Object) {
	runeMap[r.id] = obj
}
