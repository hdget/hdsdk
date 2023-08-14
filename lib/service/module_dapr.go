package service

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/dapr"
	"github.com/hdget/hdsdk/utils"
	"strings"
)

// DaprServiceModule 服务模块的方法信息
type DaprServiceModule struct {
	*BaseModule
}

func (d *DaprServiceModule) NewHandler(handler common.ServiceInvocationHandler) common.ServiceInvocationHandler {
	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 进入函数
		for _, fn := range d.EnterFunctions {
			_ = fn(ctx)
		}

		// 挂载defer函数
		defer func() {
			d.DeferFunc(d.App)
		}()

		response, err := handler(ctx, event)
		if err != nil {
			// 执行错误处理函数
			for _, fn := range d.FailedFunctions {
				_ = fn(ctx)
			}

			hdsdk.Logger.Error("handle", "name", d.Name, "err", err, "req", utils.BytesToString(event.Data))
			return dapr.Error(err)
		}

		// 执行成功处理函数
		for _, fn := range d.SuccessFunctions {
			_ = fn(ctx)
		}
		return dapr.Success(event, response)
	}
}

func (d *DaprServiceModule) GetHandlers() {
	for name, handler := range utils.StructGetReceivers(d, common.ServiceInvocationHandler(nil)) {
		h, ok := handler.(common.ServiceInvocationHandler)
		if ok {
			d.Handlers[name] = d.NewHandler(h)
		}
	}
}

// getFullHandlerName 获取完整的receiver.funcName
func (d *DaprServiceModule) getFullHandlerName(method string) string {
	return strings.Join([]string{d.App, d.Name, method}, "_")
}
