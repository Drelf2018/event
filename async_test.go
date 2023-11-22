package event_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/Drelf2018/event"
)

func TestAsyncEvent(t *testing.T) {
	evt := event.Default[int]()

	evt.OnEqual("id1", func(e *event.Event, v ...int) {
		t.Logf("OnEqual:  %v\n", v)
		e.Abort()
	}, func(e *event.Event, v ...int) {
		panic(errors.New("this error will not panic"))
	})

	evt.OnRegexp(`id\d+`, func(e *event.Event, v ...int) {
		t.Logf("OnRegexp: %v\n", v)
	})

	event.Heartbeat(2, 1, func(e *event.Event, i int) {
		t.Logf("Heartbeat#%v\n", i)
		evt.Dispatch("id"+strconv.Itoa(i), i<<i, 2*i+1)
		if i == 2 {
			e.Abort()
		}
	})
}
