package configs

type Config interface {
	Fields() []func(string) Parametr
}

type Parametr struct {
	key      string
	valType  string
	asString string
}
