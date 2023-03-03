package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/utils"
	"strings"
)

// serviceModuleHandler 服务模块的方法信息
type serviceModuleHandler struct {
	name              string                          // 模块中定义的handler对应的名字
	method            string                          // dapr调用时的完整方法名,这里会将serviceModule信息填充进去
	invocationHandler common.ServiceInvocationHandler // dapr.ServiceInvocationHandler实例
}

func (sm *ServiceModule) registerHandlers(methods map[string]common.ServiceInvocationHandler) error {
	for methodName, handler := range methods {
		k := getFullHandlerName(sm.name, utils.GetFuncName(handler))
		if _, exist := sm.handlers[k]; exist {
			return fmt.Errorf("duplicate handler registered, handler: %s", k)
		}

		sm.handlers[k] = &serviceModuleHandler{
			name:              methodName,
			method:            getFullMethodName(sm.version, sm.namespace, sm.client, methodName),
			invocationHandler: wrapRecoverHandler(sm.app, handler), // 封装recover
		}
	}
	return nil
}

// getFullHandlerName 获取完整的receiver.funcName
func getFullHandlerName(moduleName, handlerName string) string {
	return strings.Join([]string{moduleName, handlerName}, "_")
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
