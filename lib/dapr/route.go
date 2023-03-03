package dapr

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
	Version       int      // version
	Endpoint      string   // endpoint
	Methods       []string // http methods
	CallerId      int64    // 第三方回调应用id
	IsRawResponse bool     // 是否返回原始消息
	IsPublic      bool     // 是否是公共方法
	Comments      []string // 备注
}

// GetRoutes 获取路由
func (sm *ServiceModule) GetRoutes(args ...string) ([]*ServiceModuleRoute, error) {
	annotations, err := sm.GetAnnotations(args...)
	if err != nil {
		return nil, err
	}

	routes := make([]*ServiceModuleRoute, 0)
	for _, handlerAnnotation := range annotations {
		if handlerAnnotation.ModuleName != sm.name {
			continue
		}

		// 获取该handler的路由注解
		ann := handlerAnnotation.Annotations[annotationRoute]
		if ann == nil {
			continue
		}

		route, err := sm.buildRoute(handlerAnnotation.HandlerName, ann, handlerAnnotation.Comments)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (sm *ServiceModule) buildRoute(handlerName string, ann *annotation, comments []string) (*ServiceModuleRoute, error) {
	k := getFullHandlerName(sm.name, handlerName)
	handler := sm.handlers[k]
	if handler == nil {
		return nil, fmt.Errorf("handler not found, handler: %s", k)
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

	route.Version = sm.version
	route.Namespace = sm.namespace
	route.App = sm.app
	route.Handler = handler.name
	route.Comments = comments
	return route, nil
}
