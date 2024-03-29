package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"strings"
)

type invocationHandler struct {
	id     string //
	method string // receiver.method对应的方法名，例如： (*aaa) GetIdHandler(), 这里GetIdHandler为method名
	alias  string // 别名
	fn     any    // 具体的调用函数
}

type HandlerMatch func(funcName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的method名
type InvocationHandler func(ctx context.Context, event *common.InvocationEvent) (any, error)

// DiscoverHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (b *invocationModuleImpl) DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) {
	matchFn := b.defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make(map[string]any)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, method := range hdutils.Reflect().MatchReceiverMethods(b.concrete, common.ServiceInvocationHandler(nil)) {
		alias, matched := matchFn(methodName)
		if !matched {
			continue
		}

		fn := b.toDaprServiceInvocationHandler(method)
		if fn == nil {
			return nil, errors.New("invalid dpar invocation handler")
		}

		handlers[alias] = fn
	}

	return handlers, nil
}

// RegisterHandlers 参数handlers为alias=>receiver.method, 保存为handler.id=>*invocationHandler
func (b *invocationModuleImpl) RegisterHandlers(handlers map[string]any) error {
	concreteModule, ok := b.concrete.(InvocationModule)
	if !ok {
		return errors.New("invalid service invocation module")
	}

	// 校验handler
	for _, fn := range handlers {
		err := concreteModule.ValidateHandler(fn)
		if err != nil {
			return err
		}
	}

	// 这里需要将alias=>receiver.method转换成内部标识的id=>invocationHandler
	for alias, fn := range handlers {
		h := newInvocationHandler(concreteModule.GetName(), alias, fn)
		if _, exist := b.handlers[h.id]; exist {
			return errDuplicateHandler
		}
		b.handlers[h.id] = h
	}

	registerInvocationModule(b.Name, concreteModule)
	return nil
}

// GetHandlers 将map[string]*invocationHandler转换成map[string]common.serviceHandler
func (b *invocationModuleImpl) GetHandlers() map[string]any {
	handlers := make(map[string]any)
	// h为*invocationHandler
	for _, h := range b.handlers {
		// daprMethodName = v2:xxx:alias
		daprMethodName := GetServiceInvocationName(b.ModuleVersion, b.Module, h.alias)
		daprMethod := b.toDaprServiceInvocationHandler(h.fn)
		if daprMethod != nil {
			handlers[daprMethodName] = daprMethod
		}

	}
	return handlers
}

func (b *invocationModuleImpl) ValidateHandler(handler any) error {
	if hdutils.Reflect().GetFuncSignature(handler) != hdutils.Reflect().GetFuncSignature(InvocationHandler(nil)) {
		return fmt.Errorf("invalid handler: %s, it should be: func(ctx context.Context, event *common.InvocationEvent) (any, error)", hdutils.Reflect().GetFuncName(handler))
	}
	return nil
}

// toDaprServiceInvocationHandler 封装handler
func (b *invocationModuleImpl) toDaprServiceInvocationHandler(method any) common.ServiceInvocationHandler {
	realHandler, ok := method.(InvocationHandler)
	// 如果不是DaprInvocationHandler, 可能为实际的函数体
	if !ok {
		realHandler, ok = method.(func(context.Context, *common.InvocationEvent) (any, error))
		if !ok {
			return nil
		}
	}

	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				hdutils.RecordErrorStack(b.App)
			}
		}()

		response, err := realHandler(ctx, event)
		if err != nil {
			hdsdk.Logger.Error("handle", "module", b.Name, "method", hdutils.Reflect().GetFuncName(method), "err", err, "req", hdutils.BytesToString(event.Data))
			return Error(err)
		}

		return Success(event, response)
	}
}

// matchHandlerSuffix 匹配方法名是否以handler结尾并将新方法名转为SnakeCase格式
func (b *invocationModuleImpl) defaultHandlerMatchFunction(methodName string) (string, bool) {
	lowerName := strings.ToLower(methodName)
	lastIndex := strings.LastIndex(lowerName, "handler")
	if lastIndex <= 0 {
		return "", false
	}
	// handler字符串长度为7, 确保handler结尾
	if lowerName[lastIndex+7:] != "" {
		return "", false
	}
	return lowerName[:lastIndex], true
}

func newInvocationHandler(moduleName, alias string, fn any) *invocationHandler {
	methodName := hdutils.Reflect().GetFuncName(fn)
	return &invocationHandler{
		id:     genHandlerId(moduleName, methodName),
		method: methodName,
		alias:  alias,
		fn:     fn,
	}
}

// genHandlerId 生成handlerId
func genHandlerId(moduleName, methodName string) string {
	return strings.Join([]string{moduleName, methodName}, "_")
}
