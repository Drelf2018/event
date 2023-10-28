package event_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/event"
)

func TestAsyncEvent(t *testing.T) {
	a := event.Default[int]()

	a.OnCommand("av1", func(e *event.Event, v ...int) {
		fmt.Printf("OnCommand: %v (%T)\n", v, v)
		e.Abort()
	}, func(e *event.Event, v ...int) {
		panic(fmt.Errorf("this error will not panic"))
	})

	a.OnRegexp(`av\d+`, func(e *event.Event, v ...int) {
		fmt.Printf("OnRegexp:  %v (%T)\n", v, v)
	})

	a.OnAll(func(e *event.Event, v ...any) {
		fmt.Printf("OnAll:     %v (%T)\n", v, v)
	})

	a.Dispatch("TestOnAll", "text")
	event.Heartbeat(2, 1, func(e *event.Event, i int) {
		fmt.Printf("\nHeartbeat#%v\n", i)
		a.Dispatch(fmt.Sprintf("av%v", i), i<<i)
		if i == 2 {
			e.Abort()
		}
	})
}
