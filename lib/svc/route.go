package svc

import (
	"encoding/json"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"strings"
)

type Route struct {
	*moduleInfo
	Handler       string   // dapr method
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

// parseRoutes 从源代码的注解中解析路由
func (b *baseInvocationModule) parseRoutes(srcPath, annotationTag string, fnParams, fnResults []string, args ...HandlerMatch) ([]*Route, error) {
	matchFn := defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := utils.AST().InspectFunction(srcPath, fnParams, fnResults, annotationTag)
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
		ann := fnInfo.Annotations[annotationTag]
		if ann == nil {
			continue
		}

		// 忽略不匹配的函数
		_, matched := matchFn(fnInfo.Function)
		if !matched {
			continue
		}

		// 通过handlerId获取注册时候的invocationHandler
		h := b.handlers[genHandlerId(b.Name, fnInfo.Function)]
		if h == nil {
			continue
		}

		route, err := b.buildRoute(h.alias, fnInfo, ann)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}

// buildRoute alias为register或者discover handler时候使用的别名
func (b *baseInvocationModule) buildRoute(alias string, fnInfo *utils.AstFunction, ann *utils.AstAnnotation) (*Route, error) {
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
		moduleInfo:    b.moduleInfo,
		Handler:       alias,
		Endpoint:      routeAnnotation.Endpoint,
		Methods:       routeAnnotation.Methods,
		CallerId:      routeAnnotation.CallerId,
		IsRawResponse: routeAnnotation.IsRawResponse,
		IsPublic:      routeAnnotation.IsPublic,
		Comments:      fnInfo.PlainComments,
	}, nil
}
