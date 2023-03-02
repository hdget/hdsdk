package dapr

import (
	"encoding/json"
	"fmt"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"strings"
)

type serviceModuleRoute struct {
	App           string
	Handler       string
	Namespace     string
	Version       int
	Endpoint      string
	HttpMethods   []string
	CallerId      int64
	IsRawResponse bool
	IsPublic      bool
	Comments      []string
}

// GetRoutes 获取路由
func (sm *ServiceModule) GetRoutes(args ...string) ([]*serviceModuleRoute, error) {
	annotations, err := sm.GetAnnotations(args...)
	if err != nil {
		return nil, err
	}

	routes := make([]*serviceModuleRoute, 0)
	for _, handlerAnnotation := range annotations {
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

func (sm *ServiceModule) buildRoute(handlerName string, ann *annotation, comments []string) (*serviceModuleRoute, error) {
	handler := sm.handlers[handlerName]
	if handler == nil {
		return nil, fmt.Errorf("handler not found, handler: %s", handlerName)
	}

	// 尝试将注解后的值进行jsonUnmarshal
	var route *serviceModuleRoute
	v := strings.TrimSpace(ann.Value)
	if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(utils.StringToBytes(ann.Value), &route)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	} else {
		route = &serviceModuleRoute{}
	}

	route.Version = sm.version
	route.Namespace = sm.namespace
	route.App = sm.app
	route.Handler = handler.name
	route.Comments = comments
	return route, nil
}
