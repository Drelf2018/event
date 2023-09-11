package event

import (
	"time"

	"github.com/Drelf2020/utils"
)

func WithData[T any](handles ...func(*Event, T)) func(*Event) {
	return func(e *Event) {
		var t T
		utils.PanicErr(e.Data(&t))
		for _, h := range handles {
			if e.aborted {
				break
			}
			h(e, t)
		}
	}
}

func OnlyData[T any](handle func(T)) func(*Event) {
	return func(e *Event) {
		var t T
		utils.PanicErr(e.Data(&t))
		handle(t)
	}
}

func Heartbeat(initdead, keepalive float64, f func(*Event)) {
	time.Sleep(time.Duration(initdead) * time.Second)

	ticker := time.NewTicker(time.Duration(keepalive) * time.Second)
	defer ticker.Stop()

	e := eventPool.Get("Heartbeat", 0)
	e.ch = make(chan any)
	defer eventPool.Put(e)

	do := func() {
		f(e)
		e.data = e.data.(int) + 1
	}

	go do()
	for {
		select {
		case <-ticker.C:
			go do()
		case <-e.ch:
			close(e.ch)
			return
		}
	}
}
