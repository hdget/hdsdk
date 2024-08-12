package dapr

import (
	"encoding/json"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdutils/ast"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
	"strings"
)

type routeAnnotation struct {
	Endpoint      string   `json:"endpoint"`      // endpoint
	Methods       []string `json:"methods"`       // http methods
	Origin        string   `json:"origin"`        // 请求来源
	IsRawResponse bool     `json:"isRawResponse"` // 是否返回原始消息
	IsPublic      bool     `json:"isPublic"`      // 是否是公共路由
	Permissions   []string `json:"permissions"`   // 对应的权限列表
	HandlerAlias  string
	Comments      []string
}

const (
	annotationPrefix = "@hd." // 注解的前缀
	annotationRoute  = annotationPrefix + "route"
)

var (
	handlerParamSignatures  = []string{"context.Context", "*common.InvocationEvent"}
	handlerResultSignatures = []string{"any", "error"}
)

// GetRouteAnnotations 从源代码的注解中解析路由注解
func (m *invocationModuleImpl) GetRouteAnnotations(srcPath string, args ...HandlerNameMatcher) ([]*routeAnnotation, error) {
	matchFn := m.defaultHandlerNameMatcher
	if len(args) > 0 {
		matchFn = args[0]
	}

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := ast.InspectFunction(srcPath, handlerParamSignatures, handlerResultSignatures, annotationPrefix)
	if err != nil {
		return nil, err
	}

	routeAnnotations := make([]*routeAnnotation, 0)
	for _, fnInfo := range funcInfos {
		mInfo, err := parseModuleInfo(fnInfo.Receiver)
		if err != nil {
			return nil, err
		}

		// 忽略掉不是本模块的备注
		if mInfo.StructName != m.moduler.GetModuleInfo().StructName {
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

		foundIndex := pie.FindFirstUsing(m.handlers, func(h invocationHandler) bool {
			return h.GetName() == fnInfo.Function
		})
		if foundIndex == -1 {
			continue
		}

		h := m.handlers[foundIndex]
		item, err := m.parseRouteAnnotation(h.GetAlias(), fnInfo, ann)
		if err != nil {
			return nil, err
		}

		routeAnnotations = append(routeAnnotations, item)
	}

	return routeAnnotations, nil
}

// parseRouteAnnotation handlerName为register或者discover handler时候使用的handlerAlias
func (m *invocationModuleImpl) parseRouteAnnotation(handlerAlias string, fnInfo *ast.Function, ann *ast.Annotation) (*routeAnnotation, error) {
	// 设置初始值
	result := &routeAnnotation{
		Endpoint:      "",
		Methods:       []string{"GET"},
		Origin:        "",
		IsRawResponse: false,
		IsPublic:      false,
		Permissions:   []string{},
		HandlerAlias:  handlerAlias,
		Comments:      fnInfo.PlainComments,
	}

	// 尝试将注解后的值进行jsonUnmarshal
	if strings.HasPrefix(ann.Value, "{") && strings.HasSuffix(ann.Value, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(convert.StringToBytes(ann.Value), &result)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	}

	// 处理特殊情况, 设置缺省值
	if len(result.Methods) == 0 {
		result.Methods = []string{"GET"}
	}

	return result, nil
}
