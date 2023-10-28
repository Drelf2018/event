package event

import (
	"fmt"
	"time"

	"github.com/Drelf2018/TypeGo/Pool"
	"github.com/Drelf2020/utils"
)

var eventPool = Pool.New(&Event{})

type Event struct {
	cmd     string
	env     map[string]any
	aborted bool
	ch      chan any
}

func (e *Event) New() {
	e.env = make(map[string]any)
	e.aborted = false
}

func (e *Event) Set(x ...any) {
	e.cmd = fmt.Sprintf("%v", x[0])
}

func (e *Event) Reset() {
	for k := range e.env {
		delete(e.env, k)
	}
	e.aborted = false
}

func (e *Event) Cmd() string {
	return e.cmd
}

func (e Event) String() string {
	return fmt.Sprintf("Event(%v, %v)", e.cmd, e.env)
}

func (e *Event) Store(name string, value any) {
	e.env[name] = value
}

func (e *Event) Get(x any, name string, none any) error {
	if y, ok := e.env[name]; ok {
		return utils.CopyAny(x, y)
	}
	return utils.CopyAny(x, none)
}

func (e *Event) Abort() {
	e.aborted = true
	if e.ch != nil {
		e.ch <- struct{}{}
	}
}

func Heartbeat(initdead, keepalive float64, f func(e *Event, count int)) {
	time.Sleep(time.Duration(1000*initdead) * time.Millisecond)

	ticker := time.NewTicker(time.Duration(1000*keepalive) * time.Millisecond)
	defer ticker.Stop()

	count := 0
	e := eventPool.Get("Heartbeat")
	e.ch = make(chan any)
	defer eventPool.Put(e)

	go f(e, count)
	for {
		select {
		case <-ticker.C:
			count++
			go f(e, count)
		case <-e.ch:
			close(e.ch)
			return
		}
	}
}
