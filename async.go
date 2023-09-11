package event

import (
	"regexp"

	"github.com/Drelf2018/asyncio"
)

const (
	MIN int    = -2147483648
	ALL string = "__ALL__"
)

type AsyncEvent[K comparable] struct {
	all    K
	models map[string]model[K]
}

func (a *AsyncEvent[K]) Register(name string, match func(cmd, key K) bool) {
	if match == nil {
		return
	}
	a.models[name] = model[K]{make(map[K]chains), match}
}

func (a *AsyncEvent[K]) On(name string, cmd K, handles ...func(*Event)) func() {
	if m, ok := a.models[name]; ok {
		m.chains[cmd] = append(m.chains[cmd], handles)
		l := len(m.chains[cmd]) - 1
		return func() { m.chains[cmd][l] = nil }
	} else {
		panic("You should register \"" + name + "\" first.")
	}
}

func (a *AsyncEvent[K]) OnCommand(cmd K, handles ...func(*Event)) func() {
	return a.On("command", cmd, handles...)
}

func (a *AsyncEvent[K]) OnAll(handles ...func(*Event)) func() {
	return a.OnCommand(a.all, handles...)
}

func (a *AsyncEvent[K]) Dispatch(cmd K, data any) {
	asyncio.Map(a.models, func(s string, m model[K]) { m.run(cmd, data) })
	if cmd != a.all {
		a.Dispatch(a.all, data)
	}
}

func New[K comparable](all K) AsyncEvent[K] {
	a := AsyncEvent[K]{all, make(map[string]model[K])}
	// Register "command" for OnCommand
	a.Register("command", func(cmd, key K) bool { return cmd == key })
	return a
}

type AsyncEventS struct {
	AsyncEvent[string]
}

func (a *AsyncEventS) OnRegexp(pattern string, handles ...func(*Event)) func() {
	return a.On("regexp", pattern, handles...)
}

func Default() AsyncEventS {
	a := New(ALL)
	// Register "regexp" for OnRegexp
	a.Register("regexp", func(cmd, key string) bool {
		matched, err := regexp.MatchString(key, cmd)
		return err == nil && matched
	})
	return AsyncEventS{a}
}
