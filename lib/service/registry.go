package service

import "github.com/pkg/errors"

var (
	registry = make(map[string]Module)
)

func Register(app string, module Module, args ...any) error {
	moduleName, err := InitializeModule(app, module, args...)
	if err != nil {
		return errors.Wrap(err, "register module")
	}

	// 注册handlers
	registry[moduleName] = module
	return nil
}

func GetRegistry() map[string]Module {
	return registry
}
