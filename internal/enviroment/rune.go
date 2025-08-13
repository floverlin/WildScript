package enviroment

import (
	"strconv"

	"github.com/fatih/color"
)

var runeCounter uint64 = 0
var runeMap = map[uint64]Object{}

type Rune struct {
	id uint64
}

func (r *Rune) Type() ObjectType { return RUNE_TYPE }
func (r *Rune) Inspect() string {
	return color.CyanString("rune " + strconv.FormatUint(r.id, 10))
}
func (r *Rune) Set(obj Object) {
	runeMap[r.id] = obj
}
func (r *Rune) Get() Object {
	obj, ok := runeMap[r.id]
	if !ok {
		panic("rune value is empty")
	}
	return obj
}
func NewRune() *Rune {
	r := &Rune{
		id: runeCounter,
	}
	runeCounter++
	runeMap[r.id] = &Nil{}
	return r
}
