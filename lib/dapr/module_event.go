package dapr

import (
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
)

type EventModule interface {
	moduler
	RegisterHandlers(functions map[string]EventFunction) error // 注册Handlers
	GetHandlers() []eventHandler                               // 获取handlers
	GetPubSub() string
}

type eventModuleImpl struct {
	moduler
	pubsub   string // 消息中间件名称定义在dapr配置中
	handlers []eventHandler
}

var (
	_ EventModule = (*eventModuleImpl)(nil)
)

// NewEventModule 新建事件模块会执行下列操作:
func NewEventModule(app, pubsub string, moduleObject EventModule, functions map[string]EventFunction, options ...EventModuleOption) error {
	// 首先实例化module
	module, err := AsEventModule(app, pubsub, moduleObject, options...)
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
func AsEventModule(app, pubsub string, moduleObject any, options ...EventModuleOption) (EventModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &eventModuleImpl{
		moduler: m,
		pubsub:  pubsub,
	}

	for _, option := range options {
		option(moduleInstance)
	}

	// 初始化module
	err = hdutils.Reflect().StructSet(moduleObject, (*EventModule)(nil), moduleInstance)
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
	m.handlers = make([]eventHandler, 0)
	for topic, fn := range functions {
		m.handlers = append(m.handlers, m.newEventHandler(m, topic, fn))
	}
	return nil
}

func (m *eventModuleImpl) GetHandlers() []eventHandler {
	return m.handlers
}

func (m *eventModuleImpl) GetPubSub() string {
	return m.pubsub
}

func (m *eventModuleImpl) newEventHandler(module EventModule, topic string, fn EventFunction) eventHandler {
	return &eventHandlerImpl{
		module: module,
		topic:  topic,
		fn:     fn,
	}
}
