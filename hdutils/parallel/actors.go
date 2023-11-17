package parallel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// SignalActor actor监听signal信息,
// 当进程收到其中一个信号或者parent context被canceled
func SignalActor(ctx context.Context, signals ...os.Signal) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	return func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, signals...)
			defer signal.Stop(c)
			select {
			case sig := <-c:
				return SignalError{Signal: sig}
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(error) {
			cancel()
		}
}

// SignalError 当signalActor的execute函数收到信号退出是返回该错误
type SignalError struct {
	Signal os.Signal
}

// Error 实现Error接口
func (e SignalError) Error() string {
	return fmt.Sprintf("received signal %s", e.Signal)
}
