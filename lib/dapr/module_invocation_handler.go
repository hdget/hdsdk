package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdutils/convert"
	panicUtils "github.com/hdget/hdutils/panic"
	reflectUtils "github.com/hdget/hdutils/reflect"
)

type invocationHandler interface {
	GetAlias() string
	GetName() string
	GetInvokeName() string                                                        // 调用名字
	GetInvokeFunction(logger intf.LoggerProvider) common.ServiceInvocationHandler // 具体的调用函数
}

type invocationHandlerImpl struct {
	module moduler
	// handler的别名，
	// 如果DiscoverHandlers调用, 会将函数名作为入参，matchFunction的返回值当作别名，缺省是去除Handler后缀并小写
	// 如果RegisterHandlers调用，会直接用map的key值当为别名
	handlerAlias string
	handlerName  string             // 调用函数名
	fn           InvocationFunction // 调用函数
}

type InvocationFunction func(ctx context.Context, event *common.InvocationEvent) (any, error)
type HandlerNameMatcher func(methodName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的handlerName

var (
	maxRequestLength = 120
)

func (h invocationHandlerImpl) GetAlias() string {
	return h.handlerAlias
}

func (h invocationHandlerImpl) GetName() string {
	return h.handlerName
}

func (h invocationHandlerImpl) GetInvokeName() string {
	return Api().GetServiceInvocationName(h.module.GetMeta().ModuleVersion, h.module.GetMeta().ModuleName, h.handlerAlias)
}

func (h invocationHandlerImpl) GetInvokeFunction(logger intf.LoggerProvider) common.ServiceInvocationHandler {
	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				panicUtils.RecordErrorStack(h.module.GetApp())
			}
		}()

		response, err := h.fn(ctx, event)
		if err != nil {
			req := []rune(convert.BytesToString(event.Data))
			if len(req) > maxRequestLength {
				req = append(req[:maxRequestLength], []rune("...")...)
			}
			logger.Error("service invoke", "module", h.module.GetMeta().StructName, "handler", reflectUtils.GetFuncName(h.fn), "err", err, "req", req)
			return Error(err)
		}

		return Success(event, response)
	}
}
