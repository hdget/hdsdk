package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/utils"
)

// serviceModuleHandler 服务模块的方法信息
type serviceModuleHandler struct {
	name              string                          // 模块中定义的handler对应的名字
	method            string                          // dapr调用时的完整方法名,这里会将serviceModule信息填充进去
	invocationHandler common.ServiceInvocationHandler // dapr.ServiceInvocationHandler实例
}

func (sm *ServiceModule) registerHandlers(methods map[string]common.ServiceInvocationHandler) {
	for methodName, handler := range methods {
		sm.handlers[utils.GetFuncName(handler)] = &serviceModuleHandler{
			name:              methodName,
			method:            getFullMethodName(sm.version, sm.namespace, sm.client, methodName),
			invocationHandler: wrapRecoverHandler(sm.app, handler), // 封装recover
		}
	}
}

// wrapRecoverHandler 将panic recover处理逻辑封装进去
func wrapRecoverHandler(app string, handler common.ServiceInvocationHandler) common.ServiceInvocationHandler {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		defer func() {
			if r := recover(); r != nil {
				utils.RecordErrorStack(app)
			}
		}()
		return handler(ctx, in)
	}
}
