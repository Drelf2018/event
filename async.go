package event

import (
	"fmt"

	"github.com/Drelf2018/asyncio"
)

const (
	EQUAL  string = "equal"
	REGEXP string = "regexp"
)

type Rule[K comparable] interface {
	Add(cmd K, handles any) (int, error)
	Del(cmd K, index int)
	Emit(cmd K, data ...any)
}

type handle[V any] func(e *Event, v ...V)

type AsyncEvent[K comparable, V any] struct {
	rules map[string]Rule[K]
}

func (a *AsyncEvent[K, V]) Register(name string, rule Rule[K]) *AsyncEvent[K, V] {
	a.rules[name] = rule
	return a
}

func (a *AsyncEvent[K, V]) Rule(name string) Rule[K] {
	if r, ok := a.rules[name]; ok {
		return r
	}
	panic(fmt.Errorf("rule \"%v\" not found", name))
}

func (a *AsyncEvent[K, V]) Add(name string, cmd K, handles ...any) (int, error) {
	return a.Rule(name).Add(cmd, handles)
}

func (a *AsyncEvent[K, V]) On(name string, cmd K, handles ...handle[V]) (int, error) {
	return a.Rule(name).Add(cmd, handles)
}

func (a *AsyncEvent[K, V]) OnEqual(cmd K, handles ...handle[V]) (int, error) {
	return a.On(EQUAL, cmd, handles...)
}

func (a *AsyncEvent[K, V]) Dispatch(cmd K, data ...any) {
	asyncio.WaitGroup(len(a.rules), func(done func()) {
		for k := range a.rules {
			go func(k string) {
				defer done()
				a.rules[k].Emit(cmd, data...)
			}(k)
		}
	})
}

func New[K comparable, V any]() AsyncEvent[K, V] {
	return AsyncEvent[K, V]{
		rules: map[string]Rule[K]{EQUAL: make(Equal[K, V])},
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
	a.Register(REGEXP, NewRegexp[V]())
	return
}
