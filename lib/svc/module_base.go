package svc

import "errors"

type baseModule struct {
	self     any // 实际module
	App      string
	Version  int    // 版本号
	Name     string // 服务模块名
	handlers map[string]any
}

func newBaseModule(app string, svcHolder any) (*baseModule, error) {
	moduleName, version, err := getModuleNameAndVersion(svcHolder)
	if err != nil {
		return nil, err
	}

	return &baseModule{
		App:     app,
		Version: version,
		Name:    moduleName,
		self:    svcHolder,
	}, nil
}

func (b *baseModule) GetName() string {
	return b.Name
}

func (m *DaprModule) RegisterHandlers(handlers map[string]any) error {
	selfModule, ok := m.self.(ServiceInvocationModule)
	if !ok {
		return errors.New("invalid service invocation module")
	}

	// 校验handler
	for handlerName, handler := range handlers {
		err := selfModule.ValidateHandler(handlerName, handler)
		if err != nil {
			return err
		}
	}

	m.handlers = handlers
	addRegistry(m.Name, selfModule)
	return nil
}
