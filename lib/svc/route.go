package svc

import (
	"encoding/json"
	"github.com/hdget/hdsdk/lib/ast"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"strings"
)

type Route struct {
	App           string   // app name
	Handler       string   // dapr method
	Namespace     string   // namespace
	Version       int      // 版本
	Endpoint      string   // endpoint
	Methods       []string // http methods
	CallerId      int64    // 第三方回调应用id
	IsRawResponse bool     // 是否返回原始消息
	IsPublic      bool     // 是否是公共方法
	Comments      []string // 备注
}

type RouteAnnotation struct {
	Endpoint      string   // endpoint
	Methods       []string // http methods
	CallerId      int64    // 第三方回调应用id
	IsRawResponse bool     // 是否返回原始消息
	IsPublic      bool     // 是否是公共方法
	Comments      []string // 备注
}

type HandlerMatch func(funcName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的method名

// ParseRoutes 从源代码的注解中解析路由
func (b *baseModule) ParseRoutes(srcPath, annotationPrefix string, fnParams, fnResults []string, args ...HandlerMatch) ([]*Route, error) {
	matchFn := defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := ast.InspectFunctionByInOut(srcPath, fnParams, fnResults, annotationPrefix)
	if err != nil {
		return nil, err
	}

	routes := make([]*Route, 0)
	for _, fnInfo := range funcInfos {
		modInfo, err := getModuleInfo(fnInfo.Receiver)
		if err != nil {
			return nil, err
		}

		// 忽略掉不是本模块的备注
		if modInfo.Namespace != b.Namespace {
			continue
		}

		// 无路由注解忽略
		ann := fnInfo.Annotations[annotationRoute]
		if ann == nil {
			continue
		}

		// 忽略不匹配的函数
		newHandlerName, matched := matchFn(fnInfo.Function)
		if !matched {
			continue
		}

		route, err := b.buildRoute(newHandlerName, fnInfo, ann)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (b *baseModule) buildRoute(handlerName string, fnInfo *ast.FunctionInfo, ann *ast.Annotation) (*Route, error) {
	// 尝试将注解后的值进行jsonUnmarshal
	var routeAnnotation *RouteAnnotation
	if strings.HasPrefix(ann.Value, "{") && strings.HasSuffix(ann.Value, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(utils.StringToBytes(ann.Value), &routeAnnotation)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	} else {
		routeAnnotation = &RouteAnnotation{}
	}

	return &Route{
		App:           b.App,
		Namespace:     b.Namespace,
		Version:       b.Version,
		Handler:       handlerName,
		Endpoint:      routeAnnotation.Endpoint,
		Methods:       routeAnnotation.Methods,
		CallerId:      routeAnnotation.CallerId,
		IsRawResponse: routeAnnotation.IsRawResponse,
		IsPublic:      routeAnnotation.IsPublic,
		Comments:      fnInfo.PlainComments,
	}, nil
}

// matchHandlerSuffix 匹配方法名是否以handler结尾并将新方法名转为SnakeCase格式
func defaultHandlerMatchFunction(methodName string) (string, bool) {
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
