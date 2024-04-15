package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"strings"
)

type InvocationModule interface {
	GetInfo() *ModuleInfo                                                                       // 获取模块基本信息
	GetApp() string                                                                             // 获取APP
	DiscoverHandlers(args ...HandlerNameMatcher) ([]InvocationHandler, error)                   // 通过反射发现Handlers
	RegisterHandlers(functions map[string]InvocationFunction) error                             // 注册Handlers
	GetHandlers() []InvocationHandler                                                           // 获取手动注册的handlers
	GetRouteAnnotations(srcPath string, args ...HandlerNameMatcher) ([]*RouteAnnotation, error) // 从源代码获取路由注解
}

type invocationModuleImpl struct {
	*ModuleInfo
	concrete any // 实际module
	App      string
	handlers []InvocationHandler
}

var (
	handlerNameSuffix    = "handler"
	errInvalidModuleName = errors.New("invalid module name, it should be: v<number>_name, e,g: v1_abc")
	errInvalidFunction   = errors.New("invalid invocation function signature, it should be: func(context.Context, *common.InvocationEvent) (any, error)")

	_ InvocationModule = (*invocationModuleImpl)(nil)
)

func AsInvocationModule(app string, moduleObject any) (InvocationModule, error) {
	modInfo, err := parseModuleInfo(moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &invocationModuleImpl{
		ModuleInfo: modInfo,
		concrete:   moduleObject,
		App:        app,
	}

	// 初始化module
	err = hdutils.Reflect().StructSet(moduleObject, (*InvocationModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "install module: %+v", moduleInstance.GetInfo())
	}

	module, ok := moduleObject.(InvocationModule)
	if !ok {
		return nil, errors.New("invalid invocation module")
	}

	return module, nil
}

// AnnotateInvocationModule 注解模块会执行下列操作:
// 1. 实例化invocation module
// 2. 注册invocation functions
// 3. 注册module
func AnnotateInvocationModule(app string, moduleObject InvocationModule, functions map[string]InvocationFunction) error {
	// 首先实例化module
	module, err := AsInvocationModule(app, moduleObject)
	if err != nil {
		return err
	}

	// 然后注册handlers
	err = module.RegisterHandlers(functions)
	if err != nil {
		return err
	}

	// 最后注册module
	registerInvocationModule(module)

	return nil
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (m *invocationModuleImpl) RegisterHandlers(functions map[string]InvocationFunction) error {
	m.handlers = make([]InvocationHandler, 0)
	for handlerAlias, fn := range functions {
		m.handlers = append(m.handlers, newInvocationHandler(m.App, handlerAlias, m.ModuleInfo, fn))
	}
	return nil
}

// DiscoverHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (m *invocationModuleImpl) DiscoverHandlers(args ...HandlerNameMatcher) ([]InvocationHandler, error) {
	matchFn := m.defaultHandlerNameMatcher
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make([]InvocationHandler, 0)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, method := range hdutils.Reflect().MatchReceiverMethods(m.concrete, InvocationFunction(nil)) {
		handlerName, matched := matchFn(methodName)
		if !matched {
			continue
		}

		fn, err := m.toInvocationFunction(method)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, newInvocationHandler(m.App, handlerName, m.ModuleInfo, fn))
	}

	return handlers, nil
}

func (m *invocationModuleImpl) GetHandlers() []InvocationHandler {
	return m.handlers
}

func (m *invocationModuleImpl) GetInfo() *ModuleInfo {
	return m.ModuleInfo
}

func (m *invocationModuleImpl) GetApp() string {
	return m.App
}

func (m *invocationModuleImpl) toInvocationFunction(fn any) (InvocationFunction, error) {
	realFunction, ok := fn.(InvocationFunction)
	// 如果不是DaprInvocationHandler, 可能为实际的函数体
	if !ok {
		realFunction, ok = fn.(func(context.Context, *common.InvocationEvent) (any, error))
		if !ok {
			return nil, errInvalidFunction
		}
	}
	return realFunction, nil
}

// matchHandlerSuffix 匹配方法名是否以handler结尾并将新方法名转为SnakeCase格式
func (m *invocationModuleImpl) defaultHandlerNameMatcher(methodName string) (string, bool) {
	lastIndex := strings.LastIndex(strings.ToLower(methodName), strings.ToLower(handlerNameSuffix))
	if lastIndex <= 0 {
		return "", false
	}
	return methodName[:lastIndex], true
}
