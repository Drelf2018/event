# event

异步事件分发类，可组合后使用。

参数 `handles` 为不定长，事件触发时按照链式顺序逐一调用，可使用 `e.Abort()` 中断。

脱胎于原项目 [event.go](https://github.com/Drelf2018/asyncio/blob/v0.8.0/event.go)

### 使用

```go
package event_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/Drelf2018/event"
	"github.com/Drelf2020/utils"
)

func TestAsyncEvent(t *testing.T) {
	a := event.Default()

	a.OnCommand("danmaku114", event.OnlyData(func(data int) {
		fmt.Printf("data: %v(%T)\n", data, data)
	}))

	a.OnRegexp(`danmaku\d`,
		event.WithData(
			func(e *event.Event, data int) {
				if data&1 == 0 {
					e.Store("sin", math.Sin(float64(data)))
				}
			},
		),
		func(e *event.Event) {
			var num float64
			err := e.Get("sin", &num, -1.0)
			utils.PanicErr(err)
			if num == -1.0 {
				println("sin: Didn't store the value of sin(data)")
			} else {
				fmt.Printf("sin: %v(%T)\n", num, num)
			}
		},
	)

	a.OnAll(
		func(e *event.Event) { fmt.Printf("%v\n", e) },
		func(e *event.Event) { e.Abort() },
		func(e *event.Event) { fmt.Println("Why still running!?") },
	)

	event.Heartbeat(1, 3, event.WithData(func(e *event.Event, count int) {
		println()
		a.Dispatch("danmaku114", count)
		if count == 2 {
			e.Abort()
		}
	}))
}
```

#### 控制台

```
data: 0(int)
sin: 0(float64)
Event(__ALL__, 0, map[])

data: 1(int)
sin: Didn't store the value of sin(data)
Event(__ALL__, 1, map[])

sin: 0.9092974268256816(float64)
data: 2(int)
Event(__ALL__, 2, map[])
PASS
ok      github.com/Drelf2018/event      7.049s
```