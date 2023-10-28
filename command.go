package event

import (
	"errors"

	"github.com/Drelf2018/asyncio"
	"github.com/Drelf2020/utils"
)

var ErrHandles = errors.New("wrong type of handles")

type Command[K comparable, V any] map[K][][]handle[V]

func (c Command[K, V]) Add(cmd K, handles any) (int, error) {
	if handles, ok := handles.([]handle[V]); ok {
		c[cmd] = append(c[cmd], handles)
		return len(c[cmd]) - 1, nil
	}
	return 0, ErrHandles
}

func (c Command[K, V]) Del(cmd K, index int) {
	c[cmd][index] = nil
}

func emit[K comparable, V any](handles [][]handle[V], cmd K, data ...any) {
	v := make([]V, 0, len(data))
	if utils.Try(func() {
		for _, d := range data {
			v = append(v, d.(V))
		}
	}) != nil {
		return
	}
	asyncio.ForEach(handles, func(h []handle[V]) {
		e := eventPool.Get(cmd)
		for i, l := 0, len(h); i < l && !e.aborted; i++ {
			h[i](e, v...)
		}
	})
}

func (c Command[K, V]) Emit(cmd K, data ...any) {
	if h, ok := c[cmd]; ok {
		emit(h, cmd, data...)
	}
}

type All[K comparable] [][]handle[any]

func (a *All[K]) Add(cmd K, handles any) (int, error) {
	if handles, ok := handles.([]handle[any]); ok {
		*a = append(*a, handles)
		return len(*a) - 1, nil
	}
	return 0, ErrHandles
}
func (a *All[K]) Del(cmd K, index int) {
	(*a)[index] = nil
}

func (a *All[K]) Emit(cmd K, data ...any) {
	emit(*a, cmd, data...)
}
