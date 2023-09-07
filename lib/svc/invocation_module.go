package svc

import (
	"github.com/pkg/errors"
	"strings"
)

type InvocationModule interface {
	GetName() string
	GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error)
	DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) // 通过反射发现Handlers
	RegisterHandlers(handlers map[string]any) error
	GetHandlers() map[string]any // 获取手动注册的handlers
	ValidateHandler(handler any) error
	GetPermissions(srcPath string, args ...HandlerMatch) ([]*Permission, error)
}

type moduleInfo struct {
	Name    string // 结构名
	Module  string // 模块名， 结构名后去掉v<版本号>_的部分
	Version int    // 版本号
}

type HandlerMatch func(funcName string) (string, bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的method名

var (
	errInvalidModuleName = errors.New("invalid module name, it should be: v<number>_name, e,g: v1_abc")
)

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
