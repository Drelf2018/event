package event

import (
	"errors"

	"github.com/Drelf2018/asyncio"
	"github.com/Drelf2020/utils"
)

var ErrHandles = errors.New("wrong type of handles")

type Equal[K comparable, V any] map[K][][]handle[V]

func (e Equal[K, V]) Add(cmd K, handles any) (int, error) {
	if handles, ok := handles.([]handle[V]); ok {
		e[cmd] = append(e[cmd], handles)
		return len(e[cmd]) - 1, nil
	}
	return 0, ErrHandles
}

func (e Equal[K, V]) Del(cmd K, index int) {
	e[cmd][index] = nil
}

func AnySliceTo[T any](x []any) ([]T, error) {
	t := make([]T, len(x))
	return t, utils.Try(func() {
		for i, v := range x {
			t[i] = v.(T)
		}
	})
}

func Emit[K comparable, V any](handles [][]handle[V], cmd K, data ...any) {
	v, err := AnySliceTo[V](data)
	if err != nil {
		return
	}
	asyncio.ForEach(handles, func(h []handle[V]) {
		e := eventPool.Get(cmd)
		for _, fn := range h {
			fn(e, v...)
			if e.aborted {
				break
			}
		}
	})
}

func (e Equal[K, V]) Emit(cmd K, data ...any) {
	if h, ok := e[cmd]; ok {
		Emit(h, cmd, data...)
	}
}
