package idata

import (
	"os"

	"github.com/nektro/go-util/types"
)

var (
	Handlers    = map[string]func(string, string){}
	Concurrency int
	Guard       *types.Semaphore
	Log, _      = os.Create("debug.log")
)
