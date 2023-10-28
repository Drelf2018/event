package event

import (
	"regexp"

	"github.com/Drelf2018/asyncio"
)

type Regexp[V any] struct {
	Command[string, V]
}

func (r Regexp[V]) Match(cmd, key string) bool {
	matched, err := regexp.MatchString(key, cmd)
	return err == nil && matched
}

func (r Regexp[V]) Emit(cmd string, data ...any) {
	asyncio.Map(r.Command, func(key string, h [][]handle[V]) {
		if r.Match(cmd, key) {
			emit(h, cmd, data...)
		}
	})
}
