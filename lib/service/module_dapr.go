package service

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/ast"
	"github.com/hdget/hdsdk/lib/dapr"
	"github.com/hdget/hdsdk/utils"
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
func (m *DaprModule) NewHandler(handler common.ServiceInvocationHandler) common.ServiceInvocationHandler {
	return func(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
		// 挂载defer函数
		defer func() {
			if r := recover(); r != nil {
				utils.RecordErrorStack(m.App)
			}
		}()

		response, err := handler(ctx, event)
		if err != nil {
			hdsdk.Logger.Error("handle", "name", m.Name, "err", err, "req", utils.BytesToString(event.Data))
			return dapr.Error(err)
		}

		return dapr.Success(event, response)
	}
}

// GetHandlers 获取Module作为receiver的所有MethodMatchFunction匹配的方法, MethodMatchFunction生成新的方法名和判断是否匹配
func (m *DaprModule) GetHandlers(args ...MethodMatchFunction) map[string]any {
	matchFn := matchHandlerSuffix
	if len(args) > 0 {
		matchFn = args[0]
	}

	handlers := make(map[string]any)
	// 这里需要传入当前实际正在使用的服务模块，即带有common.ServiceInvocationHandler的struct实例
	for methodName, handler := range utils.StructGetReceiverMethodsByType(m.realModule, common.ServiceInvocationHandler(nil)) {
		daprSvcInvocationHandler, ok := handler.(func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error))
		if ok {
			newMethodName, matched := matchFn(methodName)
			if matched {
				// 构造完整的方法名，带上app, name, method
				fullMethodName := strings.Join([]string{m.App, m.Name, newMethodName}, "_")
				handlers[fullMethodName] = m.NewHandler(daprSvcInvocationHandler)
			}
		}
	}
	return handlers
}

// GetRoutes 获取路由
func (m *DaprModule) GetRoutes(srcPath string) ([]*Route, error) {
	handlers := m.GetHandlers()

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := ast.InspectFunctionByInOut(srcPath,
		[]any{context.Background(), &common.InvocationEvent{}},
		[]any{&common.Content{}, (error)(nil)},
		annotationPrefix)
	if err != nil {
		return nil, err
	}

	routes := make([]*Route, 0)
	for _, fnInfo := range funcInfos {
		// 忽略掉不是本模块的备注
		if fnInfo.Receiver != m.Name {
			continue
		}

		// 检查发射获取的函数是否与ast获取的函数名相同
		if _, exist := handlers[fnInfo.Function]; exist {
			return nil, fmt.Errorf("handler not found, handler: %s", fnInfo.Function)
		}

		// 获取该handler的路由注解
		ann := fnInfo.Annotations[annotationRoute]
		if ann == nil {
			continue
		}

		route, err := m.buildRoute(fnInfo, ann)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}

// matchHandlerSuffix 匹配方法名是否以handler结尾并将新方法名转为SnakeCase格式
func matchHandlerSuffix(methodName string) (string, bool) {
	lowerCaseName := strings.ToLower(methodName)
	lastIndex := strings.LastIndex(lowerCaseName, "handler")
	// 7是handler字符串长度，这里检查handler后面是否还有字符来判断是否以handler结尾
	// 为什么不用HasSuffix函数去判断是想后面在返回新的方法名的时候可以重用lastIndex，不要再去做一次字符串遍历
	if methodName[lastIndex+7:] != "" {
		return "", false
	}
	return lowerCaseName[:lastIndex], true
}
