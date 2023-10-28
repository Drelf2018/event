package event

import (
	"errors"

	"github.com/Drelf2018/asyncio"
)

const (
	COMMAND string = "command"
	REGEXP  string = "regexp"
	ALL     string = "__all__"
)

var ErrNoRule = errors.New("no rule")

type rule[K comparable] interface {
	Add(cmd K, handles any) (int, error)
	Del(cmd K, index int)
	Emit(cmd K, data ...any)
}

type handle[V any] func(e *Event, v ...V)

type AsyncEvent[K comparable, V any] struct {
	rules map[string]rule[K]
}

func (a *AsyncEvent[K, V]) Register(name string, rule rule[K]) *AsyncEvent[K, V] {
	a.rules[name] = rule
	return a
}

func (a *AsyncEvent[K, V]) Rule(name string) rule[K] {
	if r, ok := a.rules[name]; ok {
		return r
	}
	panic(ErrNoRule)
}

func (a *AsyncEvent[K, V]) Add(name string, cmd K, handles ...any) (int, error) {
	return a.Rule(name).Add(cmd, handles)
}

func (a *AsyncEvent[K, V]) On(name string, cmd K, handles ...handle[V]) (int, error) {
	return a.Rule(name).Add(cmd, handles)
}

func (a *AsyncEvent[K, V]) OnAll(handles ...handle[any]) (int, error) {
	var zero K
	return a.Rule(ALL).Add(zero, handles)
}

func (a *AsyncEvent[K, V]) OnCommand(cmd K, handles ...handle[V]) (int, error) {
	return a.On(COMMAND, cmd, handles...)
}

func (a *AsyncEvent[K, V]) Dispatch(cmd K, data ...any) {
	asyncio.Map(a.rules, func(_ string, r rule[K]) { r.Emit(cmd, data...) })
}

func New[K comparable, V any]() AsyncEvent[K, V] {
	return AsyncEvent[K, V]{
		rules: map[string]rule[K]{
			COMMAND: make(Command[K, V]),
			ALL:     new(All[K]),
		},
	}
}

type AsyncEventS[V any] struct {
	AsyncEvent[string, V]
}

func (a *AsyncEventS[V]) OnRegexp(pattern string, handles ...handle[V]) (int, error) {
	return a.On(REGEXP, pattern, handles...)
}

func Default[V any]() (a AsyncEventS[V]) {
	a.AsyncEvent = New[string, V]()
	a.Register(REGEXP, Regexp[V]{make(Command[string, V])})
	return
}
