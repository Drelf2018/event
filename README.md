# event

异步事件分发类，可组合后使用。

参数 `handles` 为不定长，事件触发时按照链式顺序逐一调用，可使用 `e.Abort()` 中断。

脱胎于原项目 [event.go](https://github.com/Drelf2018/asyncio/blob/v0.8.0/event.go)

### 使用

```go
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
```

#### 控制台

```
OnAll:     [text] ([]interface {})

Heartbeat#0
OnAll:     [0] ([]interface {})
OnRegexp:  [0] ([]int)

Heartbeat#1
OnRegexp:  [2] ([]int)
OnCommand: [2] ([]int)
OnAll:     [2] ([]interface {})

Heartbeat#2
OnAll:     [8] ([]interface {})
OnRegexp:  [8] ([]int)
```