package svc

import (
	"github.com/pkg/errors"
)

type ServiceInvocationModule interface {
	GetName() string
	GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error)
	DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) // 通过反射发现Handlers
	RegisterHandlers(handlers map[string]any) error
	GetHandlers() map[string]any // 获取手动注册的handlers
	ValidateHandler(name string, handler any) error
	//GetPermGroups(srcPath string) ([]*PermGroup, error)
}

type moduleInfo struct {
	Name      string // 模块名
	Namespace string // 命名空间， 命名空间模块名后去掉v<版本号>_的部分
	Version   int    // 版本号
}

var (
	moduleRegistry       = make(map[string]ServiceInvocationModule)
	errInvalidModuleName = errors.New("invalid module name, it should be: v<number>_name, e,g: v1_abc")
)

func GetRegistry() map[string]ServiceInvocationModule {
	return moduleRegistry
}

func addRegistry(name string, module ServiceInvocationModule) {
	moduleRegistry[name] = module
}

//func RegisterDaprModule(app string, version int, svcHolder any, options ...Option) error {
//	err := utils.StructSetComplexField(svcHolder, (*ServiceInvocationModule)(nil), NewDaprModule(app, moduleName, version, options...))
//	if err != nil {
//		return errors.Wrapf(err, "set base module for: %s ", moduleName)
//	}
//
//	module, ok := svcHolder.(ServiceInvocationModule)
//	if !ok {
//		return errors.New("invalid module")
//	}
//
//	// 将实际的struct实例保存进去
//	module.Super(svcHolder)
//
//	// 注册handlers
//	moduleRegistry[module.GetName()] = module
//
//	return nil
//}
