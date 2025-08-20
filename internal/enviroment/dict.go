package enviroment

import (
	"fmt"
	"strings"
)

type Dict struct {
	strMap map[string]Object
	numMap map[float64]Object
}

func (d *Dict) Len() int {
	return len(d.numMap) + len(d.strMap)
}

func NewDict() *Dict {
	return &Dict{
		strMap: map[string]Object{},
		numMap: map[float64]Object{},
	}
}

func (d *Dict) String() string {
	var sb strings.Builder
	for key, val := range d.strMap {
		sb.WriteString(fmt.Sprintf("%s: %s, ", key, val.Inspect()))
	}
	for key, val := range d.numMap {
		sb.WriteString(fmt.Sprintf("%f: %s, ", key, val.Inspect()))
	}
	return sb.String()
}

func (d *Dict) Set(k, v Object) {
	switch k := k.(type) {
	case *Num:
		d.numMap[k.Value] = v
	case *Str:
		d.strMap[k.Value] = v
	default:
		panic("TODO")
	}
}

func (d *Dict) Get(k Object) (Object, bool) {
	switch k := k.(type) {
	case *Str:
		v, ok := d.strMap[k.Value]
		return v, ok
	case *Num:
		v, ok := d.numMap[k.Value]
		return v, ok
	default:
		panic("TODO")
	}
}
