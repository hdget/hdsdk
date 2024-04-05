package hddapr

import (
	"encoding/json"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"strings"
)

type RouteAnnotation struct {
	*moduleInfo
	Handler       string   // dapr fnName
	Endpoint      string   // endpoint
	HttpMethods   []string // http methods
	Origin        string   // 请求来源
	IsRawResponse bool     // 是否返回原始消息
	IsPublic      bool     // 是否是公共路由
	Permissions   []string // 所属权限列表
	Comments      []string // 备注
}

type rawRouteAnnotation struct {
	Endpoint      string   `json:"endpoint"`      // endpoint
	Methods       []string `json:"methods"`       // http methods
	Origin        string   `json:"origin"`        // 请求来源
	IsRawResponse bool     `json:"isRawResponse"` // 是否返回原始消息
	IsPublic      bool     `json:"isPublic"`      // 是否是公共路由
	Permissions   []string `json:"permissions"`   // 所属权限列表
}

// 注解的前缀
const annotationPrefix = "@hd."
const annotationRoute = annotationPrefix + "route"

// GetRouteAnnotations 获取路由注解
func (m *invocationModuleImpl) GetRouteAnnotations(srcPath string, args ...HandlerNameMatcher) ([]*RouteAnnotation, error) {
	return m.parseRouteAnnotations(
		srcPath,
		annotationPrefix,
		[]string{"context.Context", "*common.InvocationEvent"},
		[]string{"any", "error"},
		args...,
	)
}

// parseRouteAnnotations 从源代码的注解中解析路由注解
func (m *invocationModuleImpl) parseRouteAnnotations(srcPath, annotationPrefix string, fnParams, fnResults []string, args ...HandlerNameMatcher) ([]*RouteAnnotation, error) {
	matchFn := m.defaultHandlerNameMatcher
	if len(args) > 0 {
		matchFn = args[0]
	}

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := hdutils.AST().InspectFunction(srcPath, fnParams, fnResults, annotationPrefix)
	if err != nil {
		return nil, err
	}

	routeAnnotations := make([]*RouteAnnotation, 0)
	for _, fnInfo := range funcInfos {
		modInfo, err := toModuleInfo(fnInfo.Receiver)
		if err != nil {
			return nil, err
		}

		// 忽略掉不是本模块的备注
		if modInfo.Module != m.Module {
			continue
		}

		// 无路由注解忽略
		ann := fnInfo.Annotations[annotationRoute]
		if ann == nil {
			continue
		}

		// 忽略不匹配的函数
		_, matched := matchFn(fnInfo.Function)
		if !matched {
			continue
		}

		foundIndex := pie.FindFirstUsing(m.handlers, func(h InvocationHandler) bool {
			return h.GetName() == fnInfo.Function
		})
		if foundIndex == -1 {
			continue
		}

		h := m.handlers[foundIndex]
		routeAnnotation, err := m.toRouteAnnotation(h.GetAlias(), fnInfo, ann)
		if err != nil {
			return nil, err
		}

		routeAnnotations = append(routeAnnotations, routeAnnotation)
	}

	return routeAnnotations, nil
}

// toRouteAnnotation handlerName为register或者discover handler时候使用的handlerName
func (m *invocationModuleImpl) toRouteAnnotation(handlerName string, fnInfo *hdutils.AstFunction, ann *hdutils.AstAnnotation) (*RouteAnnotation, error) {
	// 设置初始值
	raw := &rawRouteAnnotation{
		Endpoint:      "",
		Methods:       []string{"GET"},
		Origin:        "",
		IsRawResponse: false,
		IsPublic:      false,
		Permissions:   []string{},
	}

	// 尝试将注解后的值进行jsonUnmarshal
	if strings.HasPrefix(ann.Value, "{") && strings.HasSuffix(ann.Value, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(hdutils.StringToBytes(ann.Value), &raw)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	}

	// 处理特殊情况, 设置缺省值
	if len(raw.Methods) == 0 {
		raw.Methods = []string{"GET"}
	}

	return &RouteAnnotation{
		moduleInfo:    m.moduleInfo,
		Handler:       handlerName,
		Endpoint:      raw.Endpoint,
		HttpMethods:   raw.Methods,
		Origin:        raw.Origin,
		IsRawResponse: raw.IsRawResponse,
		IsPublic:      raw.IsPublic,
		Permissions:   raw.Permissions,
		Comments:      fnInfo.PlainComments,
	}, nil
}
