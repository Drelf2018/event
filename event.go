package event

import (
	"fmt"

	"github.com/Drelf2018/TypeGo/Pool"
	"github.com/Drelf2020/utils"
	"golang.org/x/exp/maps"
)

var eventPool = Pool.New(&Event{})

type Event struct {
	cmd     string
	data    any
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
	e.data = x[1]
}

func (e *Event) Reset() {
	maps.Clear(e.env)
	e.aborted = false
}

func (e *Event) Cmd() string {
	return e.cmd
}

func (e Event) String() string {
	return fmt.Sprintf("Event(%v, %v, %v)", e.cmd, e.data, e.env)
}

func (e *Event) Data(x any) error {
	return utils.CopyAny(x, e.data)
}

func (e *Event) Store(name string, value any) {
	e.env[name] = value
}

func (e *Event) Get(name string, x any, _default any) error {
	if y, ok := e.env[name]; ok {
		return utils.CopyAny(x, y)
	}
	return utils.CopyAny(x, _default)
}

func (e *Event) Abort() {
	e.aborted = true
	if e.ch != nil {
		e.ch <- struct{}{}
	}
}
