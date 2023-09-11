package event

import (
	"github.com/Drelf2018/asyncio"
	"github.com/Drelf2020/utils"
)

type chain []func(*Event)

func (c chain) start(e *Event) {
	defer eventPool.Put(e)
	for i := 0; i < len(c) && !e.aborted; i++ {
		c[i](e)
	}
}

type chains []chain

func (cs chains) call(cmd, data any) {
	for _, c := range utils.NotNilSlice(cs) {
		go c.start(eventPool.Get(cmd, data))
	}
}

type model[K comparable] struct {
	chains map[K]chains
	match  func(cmd, key K) bool
}

func (m *model[K]) run(cmd K, data any) {
	asyncio.Map(m.chains, func(key K, cs chains) {
		if m.match(cmd, key) {
			cs.call(cmd, data)
		}
	})
}
