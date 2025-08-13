package enviroment

import (
	"fmt"

	"github.com/fatih/color"
)

var runeCounter uint64 = 0
var runeMap = map[string]*Rune{}
var runeObject = map[uint64]Object{}

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

func NewRune(name string) *Rune {
	if r, ok := runeMap[name]; ok {
		return r
	}
	
	r := &Rune{
		id: runeCounter,
	}
	runeMap[name] = r
	runeObject[runeCounter] = &Nil{}

	runeCounter++
	return r
}

func (r *Rune) Get() Object {
	obj, ok := runeObject[r.id]
	if !ok {
		panic("rune value is empty")
	}
	return obj
}

func (r *Rune) Set(obj Object) Object {
	runeObject[r.id] = obj
	return obj
}

func FindRune(name string) (*Rune, bool) {
	r, ok := runeMap[name]
	return r, ok
}
