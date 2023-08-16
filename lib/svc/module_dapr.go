package svc

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/dapr"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"strings"
)

// // 注解的前缀
const annotationPrefix = "@hd."
const annotationRoute = annotationPrefix + "route"

// DaprModule 服务模块的方法信息
type DaprModule struct {
	*BaseModule
}

func NewDaprModule(app, name string, version int, options ...Option) Module {
	m := &DaprModule{
		BaseModule: NewBaseModule(app, name, version),
	}

	for _, option := range options {
		option(m)
	}

	return m
}
func (m *DaprModule) NewHandler(handler any) common.ServiceInvocationHandler {
	realHandler, ok := handler.(func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error))
	if !ok {
		return nil
	}

	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				utils.RecordErrorStack(m.App)
			}
		}()

		response, err := realHandler(ctx, event)
		if err != nil {
			hdsdk.Logger.Error("handle", "name", m.Name, "err", err, "req", utils.BytesToString(event.Data))
			return dapr.Error(err)
		}

		return dapr.Success(event, response)
	}
}

// GetServiceHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (m *DaprModule) GetServiceHandlers(args ...HandlerMatch) (map[string]any, error) {
	matchFn := defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make(map[string]any)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, handler := range utils.StructGetReceiverMethodsByType(m.realModule, common.ServiceInvocationHandler(nil)) {
		newHandler := m.NewHandler(handler)
		if newHandler == nil {
			return nil, errors.New("invalid common.ServiceInvocationHandler")
		}

		newHandlerName, matched := matchFn(methodName)
		if matched {
			daprServiceMethodName := dapr.GetServiceMethodName(m.Version, m.Name, newHandlerName)
			handlers[daprServiceMethodName] = newHandler
		}
	}

	return handlers, nil
}

// FormatRouteHandlerName 构造路由方法名
func (m *DaprModule) FormatRouteHandlerName(origHandlerName string) string {
	lowerName := strings.ToLower(origHandlerName)
	lastIndex := strings.LastIndex(lowerName, "handler")
	return strings.ToLower(lowerName[:lastIndex])
}

// GetRoutes 获取路由
func (m *DaprModule) GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error) {
	return m.ParseRoutes(
		srcPath,
		annotationPrefix,
		[]string{"context.Context", "*common.InvocationEvent"},
		[]string{"*common.Content", "error"},
		args...,
	)
}
