package svc

import "errors"

type baseInvocationModule struct {
	*moduleInfo
	concrete any // 实际module
	App      string
	handlers map[string]*invocationHandler
}

var (
	errDuplicateHandler = errors.New("duplicate handler")
)

func newBaseInvocationModule(app string, svcHolder any) (*baseInvocationModule, error) {
	modInfo, err := getModuleInfo(svcHolder)
	if err != nil {
		return nil, err
	}

	return &baseInvocationModule{
		moduleInfo: modInfo,
		concrete:   svcHolder,
		App:        app,
		handlers:   make(map[string]*invocationHandler),
	}, nil
}

func (b *baseInvocationModule) GetName() string {
	return b.Name
}

// RegisterHandlers 参数handlers为alias=>receiver.method, 保存为handler.id=>*invocationHandler
func (b *baseInvocationModule) RegisterHandlers(handlers map[string]any) error {
	concreteModule, ok := b.concrete.(InvocationModule)
	if !ok {
		return errors.New("invalid service invocation module")
	}

	// 校验handler
	for _, fn := range handlers {
		err := concreteModule.ValidateHandler(fn)
		if err != nil {
			return err
		}
	}

	// 这里需要将alias=>receiver.method转换成内部标识的id=>invocationHandler
	for alias, fn := range handlers {
		h := newInvocationHandler(concreteModule.GetName(), alias, fn)
		if _, exist := b.handlers[h.id]; exist {
			return errDuplicateHandler
		}
		b.handlers[h.id] = h
	}

	addInvocationModule(b.Name, concreteModule)
	return nil
}
