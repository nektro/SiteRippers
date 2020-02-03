package idata

import (
	"github.com/nektro/go-util/types"
)

var (
	Handlers    = map[string]func(string, string){}
	Concurrency int
	Guard       *types.Semaphore
)
