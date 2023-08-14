package service

import (
	"encoding/json"
	"fmt"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"strings"
)

type ServiceModuleRoute struct {
	App           string   // app name
	Handler       string   // dapr method
	Namespace     string   // namespace
	Client        string   // 客户端
	Version       int      // version
	Endpoint      string   // endpoint
	Methods       []string // http methods
	CallerId      int64    // 第三方回调应用id
	IsRawResponse bool     // 是否返回原始消息
	IsPublic      bool     // 是否是公共方法
	Comments      []string // 备注
}

// GetRoutes 获取路由
func (b *BaseModule) GetRoutes(args ...string) ([]*ServiceModuleRoute, error) {
	annotations, err := b.GetAnnotations(args...)
	if err != nil {
		return nil, err
	}

	routes := make([]*ServiceModuleRoute, 0)
	for _, handlerAnnotation := range annotations {
		// 忽略掉不是本模块的备注
		if handlerAnnotation.ModuleName != b.Name {
			continue
		}

		// 获取该handler的路由注解
		ann := handlerAnnotation.Annotations[annotationRoute]
		if ann == nil {
			continue
		}

		route, err := b.buildRoute(handlerAnnotation.HandlerName, ann, handlerAnnotation.Comments)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (b *BaseModule) buildRoute(handlerName string, ann *annotation, comments []string) (*ServiceModuleRoute, error) {
	handler := b.Handlers[handlerName]
	if handler == nil {
		return nil, fmt.Errorf("handler not found, handler: %s", handlerName)
	}

	// 尝试将注解后的值进行jsonUnmarshal
	var route *ServiceModuleRoute
	v := strings.TrimSpace(ann.Value)
	if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(utils.StringToBytes(ann.Value), &route)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	} else {
		route = &ServiceModuleRoute{}
	}

	route.Version = b.Version
	route.Namespace = b.Namespace
	route.App = b.App
	route.Handler = handlerName
	route.Comments = comments
	return route, nil
}
