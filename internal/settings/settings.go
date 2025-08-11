package settings

var Global *Settings

type Settings struct {
	Debug bool
}

func init() {
	Global = &Settings{}
}
