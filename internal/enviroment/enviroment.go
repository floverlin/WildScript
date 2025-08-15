package enviroment

type Single struct {
	Nil   Nil
	True  Bool
	False Bool
}

type Enviroment struct {
	store  map[string]Object
	single *Single
	outer  *Enviroment
}

func New(outer *Enviroment) *Enviroment {
	var e *Enviroment
	if outer != nil {
		e = &Enviroment{
			store: make(map[string]Object),
			outer: outer,
		}
	} else {
		e = &Enviroment{
			store: make(map[string]Object),
			single: &Single{
				Nil:   Nil{},
				True:  Bool{Value: true},
				False: Bool{Value: false},
			},
		}
	}
	
	loadBuiltin(e)
	return e
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Enviroment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Enviroment) GetOuter(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetOuter(name)
	}
	return obj, ok
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

func (e *Enviroment) Single() *Single {
	if e.single != nil {
		return e.single
	}
	if e.outer != nil {
		return e.outer.Single()
	}
	panic("no single and outer")
}
