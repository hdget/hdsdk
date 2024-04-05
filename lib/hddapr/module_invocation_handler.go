package hddapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdutils"
)

type InvocationHandler interface {
	GetName() string
	GetInvokeName() string                              // 调用名字
	GetInvokeFunction() common.ServiceInvocationHandler // 具体的调用函数
}

type invocationHandlerImpl struct {
	app         string
	moduleInfo  *moduleInfo
	handlerName string             // handler的名字，例如：(*aaa) GetIdHandler(), name为GetId
	fn          InvocationFunction // 调用函数
}

type InvocationFunction func(ctx context.Context, event *common.InvocationEvent) (any, error)
type HandlerNameMatcher func(methodName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的handlerName

var (
	maxRequestLength = 120
)

func newInvocationHandler(app, handlerName string, moduleInfo *moduleInfo, fn InvocationFunction) InvocationHandler {
	return &invocationHandlerImpl{
		app:         app,
		handlerName: handlerName,
		moduleInfo:  moduleInfo,
		fn:          fn,
	}
}

func (h invocationHandlerImpl) GetName() string {
	return h.handlerName
}

func (h invocationHandlerImpl) GetInvokeName() string {
	return NewApiServiceInvocation().GetServiceInvocationName(h.moduleInfo.ModuleVersion, h.moduleInfo.Name, h.handlerName)
}

func (h invocationHandlerImpl) GetInvokeFunction() common.ServiceInvocationHandler {
	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				hdutils.RecordErrorStack(h.app)
			}
		}()

		response, err := h.fn(ctx, event)
		if err != nil {
			req := []rune(hdutils.BytesToString(event.Data))
			if len(req) > maxRequestLength {
				req = append(req[:maxRequestLength], []rune("...")...)
			}
			hdsdk.Logger.Error("handle", "module", h.moduleInfo.Name, "fnName", hdutils.Reflect().GetFuncName(h.fn), "err", err, "req", req)
			return Error(err)
		}

		return Success(event, response)
	}
}
