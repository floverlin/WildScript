package enviroment

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	SELF_RUNE  = "self"
	PROTO_RUNE = "proto"
	CALL_RUNE  = "call"

	STR_RUNE = "str"

	IDX_RUNE = "idx"
	KEY_RUNE = "key"
	VAL_RUNE = "val"
)

var runeCounter uint64 = 0
var runeMap = map[string]*Rune{}
var runeObject = map[uint64]Object{}

type Rune struct {
	ID uint64
}

func (r *Rune) Type() ObjectType { return RUNE_TYPE }
func (r *Rune) Inspect() string {
	return color.CyanString(
		fmt.Sprintf(
			"rune<%d>",
			r.ID,
		),
	)
}

func TakeRune(name string) *Rune {
	if r, ok := runeMap[name]; ok {
		return r
	}

	r := &Rune{
		ID: runeCounter,
	}
	runeMap[name] = r
	runeObject[runeCounter] = &Nil{}

	runeCounter++
	return r
}

func (r *Rune) Get() Object {
	obj, ok := runeObject[r.ID]
	if !ok {
		panic("rune value is empty")
	}
	return obj
}

func (r *Rune) Set(obj Object) Object {
	runeObject[r.ID] = obj
	return obj
}

func FindRune(name string) (*Rune, bool) {
	r, ok := runeMap[name]
	return r, ok
}
