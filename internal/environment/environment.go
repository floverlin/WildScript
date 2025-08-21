package environment

type Environment struct {
	store map[string]Object
	outer *Environment
}

func New(outer *Environment) *Environment {
	e := &Environment{
		store: make(map[string]Object),
	}
	if outer != nil {
		e.outer = outer
	} else {
		e.loadBuiltin()
	}

	return e
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) (Object, bool) {
	_, ok := e.store[name]
	if ok {
		e.store[name] = val
	} else if e.outer != nil {
		_, ok = e.outer.Set(name, val)
	}
	return val, ok
}

func (e *Environment) Create(name string, val Object) (Object, bool) {
	_, ok := e.store[name]
	if ok {
		return nil, false
	}
	e.store[name] = val
	return val, true
}
