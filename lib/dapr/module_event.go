package dapr

import (
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
)

type EventModule interface {
	Moduler
	RegisterHandlers(functions map[string]EventFunction) error // 注册Handlers
	GetHandlers() []EventHandler                               // 获取handlers
}

type eventModuleImpl struct {
	Moduler
	pubsub   string // 消息中间件名称定义在dapr配置中
	handlers []EventHandler
}

var (
	_ EventModule = (*eventModuleImpl)(nil)
)

// NewEventModule 新建事件模块会执行下列操作:
func NewEventModule(app, pubsub string, moduleObject InvocationModule, functions map[string]EventFunction) error {
	// 首先实例化module
	module, err := AsEventModule(app, pubsub, moduleObject)
	if err != nil {
		return err
	}

	// 然后注册handlers
	err = module.RegisterHandlers(functions)
	if err != nil {
		return err
	}

	// 最后注册module
	registerEventModule(module)

	return nil
}

// AsEventModule 将一个any类型的结构体转换成EventModule
//
// e,g:
//
//		type v1_test struct {
//		  InvocationModule
//		}
//
//		 v := &v1_test{}
//		 im, err := AsEventModule("app",v)
//	     if err != nil {
//	      ...
//	     }
//	     im.DiscoverHandlers()
func AsEventModule(app, pubsub string, moduleObject any) (EventModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &eventModuleImpl{
		Moduler: m,
		pubsub:  pubsub,
	}

	// 初始化module
	err = hdutils.Reflect().StructSet(moduleObject, (*InvocationModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "install module: %+v", m)
	}

	module, ok := moduleObject.(EventModule)
	if !ok {
		return nil, errors.New("invalid invocation module")
	}

	return module, nil
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (m *eventModuleImpl) RegisterHandlers(functions map[string]EventFunction) error {
	m.handlers = make([]EventHandler, 0)
	for topic, fn := range functions {
		m.handlers = append(m.handlers, m.newEventHandler(m.Moduler, m.pubsub, topic, fn))
	}
	return nil
}

func (m *eventModuleImpl) GetHandlers() []EventHandler {
	return m.handlers
}

func (m *eventModuleImpl) newEventHandler(module Moduler, pubsub, topic string, fn EventFunction) EventHandler {
	return &eventHandlerImpl{
		module: module,
		pubsub: pubsub,
		topic:  topic,
		fn:     fn,
	}
}
