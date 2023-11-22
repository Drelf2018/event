package event

import (
	"regexp"

	"github.com/Drelf2018/asyncio"
)

type Regexp[V any] struct {
	Equal[string, V]
}

func (r Regexp[V]) Match(cmd, key string) bool {
	matched, err := regexp.MatchString(key, cmd)
	return err == nil && matched
}

func (r Regexp[V]) Emit(cmd string, data ...any) {
	asyncio.Map(r.Equal, func(key string, h [][]handle[V]) {
		if r.Match(cmd, key) {
			Emit(h, cmd, data...)
		}
	})
}

func NewRegexp[V any]() Regexp[V] {
	return Regexp[V]{make(Equal[string, V])}
}
