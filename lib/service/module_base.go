package service

import (
	"encoding/json"
	"github.com/hdget/hdsdk/lib/ast"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

type BaseModule struct {
	realModule any // 实际module
	App        string
	Version    int    // 版本号
	Namespace  string // 命名空间
	Name       string // 服务模块名
}

func NewBaseModule(app, name string, version int) *BaseModule {
	return &BaseModule{
		App:     app,
		Version: version,
		Name:    name,
	}
}

func (b *BaseModule) GetName() string {
	return b.Name
}

func (b *BaseModule) Super(m any) {
	b.realModule = m
}

func (b *BaseModule) buildRoute(fnInfo *ast.FunctionInfo, ann *ast.Annotation) (*Route, error) {
	// 尝试将注解后的值进行jsonUnmarshal
	var routeAnnotation *RouteAnnotation
	// 如果定义不为空，尝试unmarshal
	err := json.Unmarshal(utils.StringToBytes(ann.Value), &routeAnnotation)
	if err != nil {
		return nil, errors.Wrapf(err, "parse route Annotation, Annotation: %s", ann.Value)
	}

	return &Route{
		App:           b.App,
		Handler:       fnInfo.Function,
		Namespace:     fnInfo.Function,
		Version:       routeAnnotation.Version,
		Endpoint:      routeAnnotation.Endpoint,
		Methods:       routeAnnotation.Methods,
		CallerId:      routeAnnotation.CallerId,
		IsRawResponse: routeAnnotation.IsRawResponse,
		IsPublic:      routeAnnotation.IsPublic,
		Comments:      routeAnnotation.Comments,
	}, nil
}
