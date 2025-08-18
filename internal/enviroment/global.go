package enviroment

type GlobalObject string

var (
	GLOBAL_NIL   = &Nil{}
	GLOBAL_TRUE  = &Bool{Value: true}
	GLOBAL_FALSE = &Bool{Value: false}
)
