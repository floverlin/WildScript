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

func New() *Enviroment {
	e := &Enviroment{
		store: make(map[string]Object),
		single: &Single{
			Nil:   Nil{},
			True:  Bool{Value: true},
			False: Bool{Value: false},
		},
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

func (e *Enviroment) Single() *Single {
	if e.single != nil {
		return e.single
	}
	if e.outer != nil {
		return e.outer.Single()
	}
	panic("no single and outer")
}
