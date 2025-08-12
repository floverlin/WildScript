package enviroment

func NewBlockEnviroment(outer *Enviroment) *Enviroment {
	e := &Enviroment{
		store: make(map[string]Object),
		outer: outer,
	}
	loadBuiltin(e)
	return e
}
