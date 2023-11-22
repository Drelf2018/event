# event

异步事件分发类，可组合后使用。

参数 `handles` 为不定长，事件触发时按照链式顺序逐一调用，可使用 `e.Abort()` 中断。

脱胎于原项目 [event.go](https://github.com/Drelf2018/asyncio/blob/v0.8.0/event.go)

### 使用

```go
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
```

#### 控制台

```
=== RUN   TestAsyncEvent
    async_test.go:26: Heartbeat#0
    async_test.go:22: OnRegexp: [0 1]
    async_test.go:26: Heartbeat#1
    async_test.go:15: OnEqual:  [2 3]
    async_test.go:22: OnRegexp: [2 3]
    async_test.go:26: Heartbeat#2
    async_test.go:22: OnRegexp: [8 5]
--- PASS: TestAsyncEvent (4.02s)
PASS
ok      github.com/Drelf2018/event      4.047s
```