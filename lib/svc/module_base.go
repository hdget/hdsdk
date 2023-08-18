package svc

import "errors"

type baseModule struct {
	*moduleInfo
	self     any // 实际module
	App      string
	handlers map[string]any
}

func newBaseModule(app string, svcHolder any) (*baseModule, error) {
	modInfo, err := getModuleInfo(svcHolder)
	if err != nil {
		return nil, err
	}

	return &baseModule{
		moduleInfo: modInfo,
		self:       svcHolder,
		App:        app,
		handlers:   nil,
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
