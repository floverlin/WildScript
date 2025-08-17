package enviroment

type Enviroment struct {
	store     map[string]Object
	runeStore map[string]Object
	outer     *Enviroment
}

func New(outer *Enviroment) *Enviroment {
	var e *Enviroment
	if outer != nil {
		e = &Enviroment{
			store:     make(map[string]Object),
			runeStore: make(map[string]Object),
			outer:     outer,
		}
	} else {
		e = &Enviroment{
			store:     make(map[string]Object),
			runeStore: make(map[string]Object),
		}
		e.loadBuiltin()
	}

	return e
}

func (e *Enviroment) GetRuneOuter(name string) (Object, bool) {
	obj, ok := e.runeStore[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetRuneOuter(name)
	}
	return obj, ok
}

func (e *Enviroment) SetRune(name string, val Object) Object {
	e.runeStore[name] = val
	return val
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Enviroment) GetOuter(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetOuter(name)
	}
	return obj, ok
}

func (e *Enviroment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Enviroment) SetOuter(name string, val Object) (Object, bool) {
	_, ok := e.store[name]
	if ok {
		e.store[name] = val
	} else if e.outer != nil {
		_, ok = e.outer.SetOuter(name, val)
	}
	return val, ok
}
