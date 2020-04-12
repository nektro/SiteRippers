package iutil

import (
	"sync"
)

var (
	mt = new(sync.Mutex)
)

type Counter struct {
	after int
	f     func()
	value int
	wg    *sync.WaitGroup
}

func NewCounter(a int, f func()) *Counter {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	return &Counter{a, f, 0, wg}
}

func (c *Counter) Add(delta int) {
	mt.Lock()
	c.after += delta
	mt.Unlock()
}

func (c *Counter) Increment() {
	mt.Lock()
	c.value++
	if c.value == c.after {
		c.f()
		c.wg.Done()
	}
	mt.Unlock()
}

func (c *Counter) Wait() {
	c.wg.Wait()
}
