package enviroment

type Single struct {
	Nil   Object
	True  Object
	False Object
}

type Environment struct {
	store  map[string]Object
	Single *Single
}

func New() *Environment {
	e := &Environment{
		store: make(map[string]Object),
		Single: &Single{
			Nil:   &Nil{},
			True:  &Bool{Value: true},
			False: &Bool{Value: false},
		},
	}
	e.loadBuiltin()
	return e
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
