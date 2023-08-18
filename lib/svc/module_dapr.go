package svc

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/dapr"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

// // 注解的前缀
const annotationPrefix = "@hd."
const annotationRoute = annotationPrefix + "route"

type Option func(module ServiceInvocationModule) ServiceInvocationModule
type DaprServiceInvocationHandler func(ctx context.Context, eventData []byte) (any, error)

// DaprModule 服务模块的方法信息
type DaprModule struct {
	*baseModule
}

func RegisterAsDaprModule(app string, svcHolder any, args ...map[string]DaprServiceInvocationHandler) error {
	module, err := NewDaprModule(app, svcHolder)
	if err != nil {
		return err
	}

	// 注册handlers
	handlers := make(map[string]any)
	if len(args) > 0 {
		for name, handler := range args[0] {
			handlers[name] = handler
		}
	} else {
		handlers, err = module.DiscoverHandlers()
		if err != nil {
			return errors.Wrap(err, "discover handlers")
		}
	}

	err = module.RegisterHandlers(handlers)
	if err != nil {
		return err
	}

	return nil
}

func NewDaprModule(app string, svcHolder any, options ...Option) (ServiceInvocationModule, error) {
	b, err := newBaseModule(app, svcHolder)
	if err != nil {
		return nil, err
	}

	m := &DaprModule{
		baseModule: b,
	}

	for _, option := range options {
		option(m)
	}

	// 将实例化的module设置入Module接口中
	err = utils.StructSetComplexField(svcHolder, (*ServiceInvocationModule)(nil), m)
	if err != nil {
		return nil, errors.Wrapf(err, "install module for: %s ", m.GetName())
	}

	module, ok := svcHolder.(ServiceInvocationModule)
	if !ok {
		return nil, errors.New("not ServiceInvocationModule")
	}

	return module, nil
}

// GetHandlers 获取手动注册的handlers
func (m *DaprModule) GetHandlers() map[string]any {
	svcInvocationHandlers := make(map[string]any)
	// 当前存在m.handlers中的为DaprHandler类型
	for handlerName, handler := range m.handlers {
		svcInvocationHandlers[handlerName] = m.toDaprServiceInvocationHandler(handlerName, handler)
	}
	return svcInvocationHandlers
}

// DiscoverHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (m *DaprModule) DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) {
	matchFn := defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make(map[string]any)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, handler := range utils.StructGetReceiverMethodsByType(m.self, common.ServiceInvocationHandler(nil)) {
		newHandlerName, matched := matchFn(methodName)
		if !matched {
			continue
		}

		svcInvocationName := dapr.GetServiceInvocationName(m.Version, m.Name, newHandlerName)
		svcInvocationHandler := m.toDaprServiceInvocationHandler(svcInvocationName, handler)
		if svcInvocationHandler == nil {
			return nil, errors.New("invalid common.ServiceInvocationHandler")
		}

		handlers[svcInvocationName] = svcInvocationHandler
	}

	return handlers, nil
}

//// FormatRouteHandlerName 构造路由方法名
//func (m *DaprModule) FormatRouteHandlerName(origHandlerName string) string {
//	lowerName := strings.ToLower(origHandlerName)
//	lastIndex := strings.LastIndex(lowerName, "handler")
//	return strings.ToLower(lowerName[:lastIndex])
//}

func (m *DaprModule) ValidateHandler(handlerName string, handler any) error {
	if _, ok := handler.(DaprServiceInvocationHandler); !ok {
		return fmt.Errorf("invalid handler: %s, it should be: func(ctx context.Context, event *common.InvocationEvent) (any, error)", handlerName)
	}
	return nil
}

// GetRoutes 获取路由
func (m *DaprModule) GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error) {
	return m.ParseRoutes(
		srcPath,
		annotationPrefix,
		[]string{"context.Context", "*common.InvocationEvent"},
		[]string{"any", "error"},
		args...,
	)
}

// toDaprServiceInvocationHandler 封装handler
func (m *DaprModule) toDaprServiceInvocationHandler(handlerName string, handler any) common.ServiceInvocationHandler {
	realHandler, ok := handler.(DaprServiceInvocationHandler)
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

		response, err := realHandler(ctx, event.Data)
		if err != nil {
			hdsdk.Logger.Error("handle", "namespace", m.Name, "handler", handlerName, "err", err, "req", utils.BytesToString(event.Data))
			return dapr.Error(err)
		}

		return dapr.Success(event, response)
	}
}
