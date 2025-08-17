package enviroment

type GlobalObject string

const (
	GLOBAL_NIL   = "nil"
	GLOBAL_TRUE  = "true"
	GLOBAL_FALSE = "false"
)

var Global = map[GlobalObject]Object{
	GLOBAL_NIL:   &Nil{},
	GLOBAL_TRUE:  &Bool{Value: true},
	GLOBAL_FALSE: &Bool{Value: false},
}
