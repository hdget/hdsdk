package dapr

import (
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
)

type InvocationModule interface {
	GetName() string
	GetRouteAnnotations(srcPath string, args ...HandlerMatch) ([]*RouteAnnotation, error)
	DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) // 通过反射发现Handlers
	RegisterHandlers(handlers map[string]any) error
	GetHandlers() map[string]any // 获取手动注册的handlers
	ValidateHandler(handler any) error
}

type moduleInfo struct {
	Name          string // 结构名, 格式: "v<模块版本号>_<模块名>"
	Module        string // 模块名
	ModuleVersion int    // 模块版本号
}

type invocationModuleImpl struct {
	*moduleInfo
	concrete any // 实际module
	App      string
	handlers map[string]*invocationHandler
}

var (
	errInvalidModuleName = errors.New("invalid module name, it should be: v<number>_name, e,g: v1_abc")
	errDuplicateHandler  = errors.New("duplicate handler")
)

// NewInvocationModule 实例化invocation module
func NewInvocationModule(app string, moduleObject any) (InvocationModule, error) {
	modInfo, err := parseModuleInfo(moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &invocationModuleImpl{
		moduleInfo: modInfo,
		concrete:   moduleObject,
		App:        app,
		handlers:   make(map[string]*invocationHandler),
	}

	// 将实例化的module设置入Module接口中
	err = hdutils.Reflect().StructSet(moduleObject, (*InvocationModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "install module for: %s ", moduleInstance.GetName())
	}

	module, ok := moduleObject.(InvocationModule)
	if !ok {
		return nil, errors.New("invalid invocation module")
	}

	return module, nil
}

// AsInvocationModule 将struct注册为DaprInvocationModule
func AsInvocationModule(app string, moduleObject any, args ...map[string]InvocationHandler) error {
	module, err := NewInvocationModule(app, moduleObject)
	if err != nil {
		return err
	}

	// 注册handlers, alias=>receiver.method
	handlers := make(map[string]any)
	if len(args) > 0 {
		for alias, fn := range args[0] {
			handlers[alias] = fn
		}
	} else {
		handlers, err = module.DiscoverHandlers()
		if err != nil {
			return errors.Wrap(err, "discover invocationHandlers")
		}
	}

	err = module.RegisterHandlers(handlers)
	if err != nil {
		return err
	}

	return nil
}

func (b *invocationModuleImpl) GetName() string {
	return b.Name
}
