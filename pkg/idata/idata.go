package idata

var (
	Handlers    = map[string]func(string, string){}
	Concurrency int
)
