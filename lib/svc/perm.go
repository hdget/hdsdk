package svc

import (
	"encoding/json"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"strings"
)

type Permission struct {
	*moduleInfo
	Name  string // app name
	Group string // permission group name
}

type PermissionAnnotation struct {
	Name  string // app name
	Group string // permission group name
}

const annotationPermission = annotationPrefix + "perm"

// parseRoutes 从源代码的注解中解析路由
func (b *baseInvocationModule) parsePermissions(srcPath, annotationPrefix string, fnParams, fnResults []string, args ...HandlerMatch) ([]*Permission, error) {
	matchFn := defaultHandlerMatchFunction
	if len(args) > 0 {
		matchFn = args[0]
	}

	// 这里需要匹配func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
	// 函数参数类型为: context.Context, *common.InvocationEvent
	// 函数返回结果为：
	funcInfos, err := utils.AST().InspectFunction(srcPath, fnParams, fnResults, annotationPrefix)
	if err != nil {
		return nil, err
	}

	perms := make([]*Permission, 0)
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
		ann := fnInfo.Annotations[annotationPermission]
		if ann == nil {
			continue
		}

		// 忽略不匹配的函数
		newHandlerName, matched := matchFn(fnInfo.Function)
		if !matched {
			continue
		}

		permission, err := b.buildPermission(newHandlerName, fnInfo, ann)
		if err != nil {
			return nil, err
		}

		perms = append(perms, permission)
	}

	return perms, nil
}

func (b *baseInvocationModule) buildPermission(handlerName string, fnInfo *utils.AstFunction, ann *utils.AstAnnotation) (*Permission, error) {
	// 尝试将注解后的值进行jsonUnmarshal
	var annotation *PermissionAnnotation
	if strings.HasPrefix(ann.Value, "{") && strings.HasSuffix(ann.Value, "}") {
		// 如果定义不为空，尝试unmarshal
		err := json.Unmarshal(utils.StringToBytes(ann.Value), &annotation)
		if err != nil {
			return nil, errors.Wrapf(err, "parse route annotation, annotation: %s", ann.Value)
		}
	} else {
		annotation = &PermissionAnnotation{}
	}

	return &Permission{
		moduleInfo: b.moduleInfo,
		Name:       annotation.Name,
		Group:      annotation.Group,
	}, nil
}
