package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	reflectUtils "github.com/hdget/hdutils/reflect"
	"github.com/pkg/errors"
	"strings"
)

type InvocationModule interface {
	moduler
	DiscoverHandlers(args ...HandlerNameMatcher) ([]invocationHandler, error)                   // 通过反射发现Handlers
	RegisterHandlers(functions map[string]InvocationFunction) error                             // 注册Handlers
	GetHandlers() []invocationHandler                                                           // 获取handlers
	GetRouteAnnotations(srcPath string, args ...HandlerNameMatcher) ([]*RouteAnnotation, error) // 从源代码获取路由注解
}

type invocationModuleImpl struct {
	moduler
	self     any // 实际module实例
	handlers []invocationHandler
}

var (
	errInvalidInvocationFunction                  = errors.New("invalid invocation function signature, it should be: func(context.Context, *common.InvocationEvent) (any, error)")
	_                            InvocationModule = (*invocationModuleImpl)(nil)
)

// NewInvocationModule 新建服务调用模块会执行下列操作:
// 1. 实例化invocation module
// 2. 注册invocation functions
// 3. 注册module
func NewInvocationModule(app string, moduleObject InvocationModule, functions map[string]InvocationFunction) error {
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

// AsInvocationModule 将一个any类型的结构体转换成InvocationModule
//
// e,g:
//
//		type v1_test struct {
//		  InvocationModule
//		}
//
//		 v := &v1_test{}
//		 im, err := AsInvocationModule("app",v)
//	     if err != nil {
//	      ...
//	     }
//	     im.DiscoverHandlers()
func AsInvocationModule(app string, moduleObject any) (InvocationModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &invocationModuleImpl{
		moduler: m,
		self:    moduleObject,
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*InvocationModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "install module: %+v", m)
	}

	module, ok := moduleObject.(InvocationModule)
	if !ok {
		return nil, errors.New("invalid invocation module")
	}

	return module, nil
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (m *invocationModuleImpl) RegisterHandlers(functions map[string]InvocationFunction) error {
	m.handlers = make([]invocationHandler, 0)
	for handlerAlias, fn := range functions {
		m.handlers = append(m.handlers, m.newInvocationHandler(m.moduler, handlerAlias, fn))
	}
	return nil
}

// DiscoverHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (m *invocationModuleImpl) DiscoverHandlers(args ...HandlerNameMatcher) ([]invocationHandler, error) {
	matchFn := m.defaultHandlerNameMatcher
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make([]invocationHandler, 0)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, method := range reflectUtils.MatchReceiverMethods(m.self, InvocationFunction(nil)) {
		handlerName, matched := matchFn(methodName)
		if !matched {
			continue
		}

		fn, err := m.toInvocationFunction(method)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, m.newInvocationHandler(m.moduler, handlerName, fn))
	}

	return handlers, nil
}

func (m *invocationModuleImpl) GetHandlers() []invocationHandler {
	return m.handlers
}

func (m *invocationModuleImpl) newInvocationHandler(module moduler, handlerAlias string, fn InvocationFunction) invocationHandler {
	return &invocationHandlerImpl{
		handlerAlias: handlerAlias,
		handlerName:  reflectUtils.GetFuncName(fn),
		module:       module,
		fn:           fn,
	}
}

func (m *invocationModuleImpl) toInvocationFunction(fn any) (InvocationFunction, error) {
	realFunction, ok := fn.(InvocationFunction)
	// 如果不是DaprInvocationHandler, 可能为实际的函数体
	if !ok {
		realFunction, ok = fn.(func(context.Context, *common.InvocationEvent) (any, error))
		if !ok {
			return nil, errInvalidInvocationFunction
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
