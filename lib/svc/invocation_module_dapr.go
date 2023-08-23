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

type Option func(module InvocationModule) InvocationModule
type DaprServiceInvocationHandler func(ctx context.Context, event *common.InvocationEvent) (any, error)

// DaprInvocationModule 服务模块的方法信息
type DaprInvocationModule struct {
	*baseInvocationModule
}

func RegisterAsDaprModule(app string, svcHolder any, args ...map[string]DaprServiceInvocationHandler) error {
	module, err := NewDaprInvocationModule(app, svcHolder)
	if err != nil {
		return err
	}

	// 注册handlers, alias=>receiver.method
	handlers := make(map[string]any)
	if len(args) > 0 {
		for alias, fn := range args[0] {
			handlers[alias] = fn
		}
	} else {
		handlers, err = module.DiscoverHandlers()
		if err != nil {
			return errors.Wrap(err, "discover invocationHandlers")
		}
	}

	err = module.RegisterHandlers(handlers)
	if err != nil {
		return err
	}

	return nil
}

func NewDaprInvocationModule(app string, svcHolder any, options ...Option) (InvocationModule, error) {
	b, err := newBaseInvocationModule(app, svcHolder)
	if err != nil {
		return nil, err
	}

	m := &DaprInvocationModule{
		baseInvocationModule: b,
	}

	for _, option := range options {
		option(m)
	}

	// 将实例化的module设置入Module接口中
	err = utils.Reflect().StructSet(svcHolder, (*InvocationModule)(nil), m)
	if err != nil {
		return nil, errors.Wrapf(err, "install module for: %s ", m.GetName())
	}

	module, ok := svcHolder.(InvocationModule)
	if !ok {
		return nil, errors.New("not InvocationModule")
	}

	return module, nil
}

// GetHandlers 将map[string]*invocationHandler转换成map[string]common.serviceHandler
func (m *DaprInvocationModule) GetHandlers() map[string]any {
	handlers := make(map[string]any)
	// h为*invocationHandler
	for _, h := range m.handlers {
		// daprMethodName = v2:xxx:alias
		daprMethodName := dapr.GetServiceInvocationName(m.Version, m.Namespace, h.alias)
		daprMethod := m.toDaprServiceInvocationHandler(h.fn)
		if daprMethod != nil {
			handlers[daprMethodName] = daprMethod
		}

	}
	return handlers
}

// DiscoverHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (m *DaprInvocationModule) DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) {
	matchFn := defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make(map[string]any)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, method := range utils.Reflect().MatchReceiverMethods(m.concrete, common.ServiceInvocationHandler(nil)) {
		alias, matched := matchFn(methodName)
		if !matched {
			continue
		}

		fn := m.toDaprServiceInvocationHandler(method)
		if fn == nil {
			return nil, errors.New("invalid common.invocationHandler")
		}

		handlers[alias] = fn
	}

	return handlers, nil
}

func (m *DaprInvocationModule) ValidateHandler(handler any) error {
	if utils.Reflect().GetFuncSignature(handler) != utils.Reflect().GetFuncSignature(DaprServiceInvocationHandler(nil)) {
		return fmt.Errorf("invalid handler: %s, it should be: func(ctx context.Context, event *common.InvocationEvent) (any, error)", utils.Reflect().GetFuncName(handler))
	}
	return nil
}

// GetRoutes 获取路由
func (m *DaprInvocationModule) GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error) {
	return m.parseRoutes(
		srcPath,
		annotationPrefix,
		[]string{"context.Context", "*common.InvocationEvent"},
		[]string{"any", "error"},
		args...,
	)
}

// GetPermissions 获取权限
func (m *DaprInvocationModule) GetPermissions(srcPath string, args ...HandlerMatch) ([]*Permission, error) {
	return m.parsePermissions(
		srcPath,
		annotationPrefix,
		[]string{"context.Context", "*common.InvocationEvent"},
		[]string{"any", "error"},
		args...,
	)
}

// toDaprServiceInvocationHandler 封装handler
func (m *DaprInvocationModule) toDaprServiceInvocationHandler(method any) common.ServiceInvocationHandler {
	realHandler, ok := method.(func(ctx context.Context, event *common.InvocationEvent) (any, error))
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
			hdsdk.Logger.Error("handle", "namespace", m.Name, "method", utils.Reflect().GetFuncName(method), "err", err, "req", utils.BytesToString(event.Data))
			return dapr.Error(err)
		}

		return dapr.Success(event, response)
	}
}
