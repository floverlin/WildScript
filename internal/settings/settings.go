package settings

var Global *Settings

type Settings struct {
	Debug  bool
	Tokens bool
}

func init() {
	Global = &Settings{}
}
