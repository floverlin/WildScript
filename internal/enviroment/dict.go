package enviroment

type Dict struct {
	strMap map[string]Object
	numMap map[float64]Object
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
