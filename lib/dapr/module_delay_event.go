package dapr

import (
	"github.com/cenkalti/backoff/v4"
	reflectUtils "github.com/hdget/hdutils/reflect"
	"github.com/pkg/errors"
	"time"
)

type DelayEventModule interface {
	moduler
	RegisterHandlers(functions map[string]DelayEventFunction) error // 注册Handlers
	GetHandlers() []delayEventHandler                               // 获取handlers
	GetAckTimeout() time.Duration
	GetBackOffPolicy() *backoff.ExponentialBackOff
}

type delayEventModuleImpl struct {
	moduler
	handlers      []delayEventHandler
	ackTimeout    time.Duration
	backoffPolicy *backoff.ExponentialBackOff
}

var (
	_ DelayEventModule = (*delayEventModuleImpl)(nil)
)

// NewDelayEventModule new delay event module
func NewDelayEventModule(app string, moduleObject DelayEventModule, functions map[string]DelayEventFunction, options ...DelayEventModuleOption) error {
	// 首先实例化module
	module, err := AsDelayEventModule(app, moduleObject, options...)
	if err != nil {
		return err
	}

	// 然后注册handlers
	err = module.RegisterHandlers(functions)
	if err != nil {
		return err
	}

	// 最后注册module
	registerDelayEventModule(module)

	return nil
}

// AsDelayEventModule 将一个any类型的结构体转换成DelayEventModule
//
// e,g:
//
//		type v1_test struct {
//		  DelayEventModule
//		}
//
//		 v := &v1_test{}
//		 im, err := AsDelayEventModule("app",v)
//	     if err != nil {
//	      ...
//	     }
//	     im.DiscoverHandlers()
func AsDelayEventModule(app string, moduleObject any, options ...DelayEventModuleOption) (DelayEventModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &delayEventModuleImpl{
		moduler:       m,
		ackTimeout:    defaultAckTimeout,
		backoffPolicy: getDefaultBackOffPolicy(),
	}

	for _, option := range options {
		option(moduleInstance)
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*DelayEventModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "install module: %+v", m)
	}

	module, ok := moduleObject.(DelayEventModule)
	if !ok {
		return nil, errors.New("invalid delay event module")
	}

	return module, nil
}

// RegisterHandlers 参数handlers为alias=>receiver.fnName, 保存为handler.id=>*invocationHandler
func (m *delayEventModuleImpl) RegisterHandlers(functions map[string]DelayEventFunction) error {
	m.handlers = make([]delayEventHandler, 0)
	for topic, fn := range functions {
		m.handlers = append(m.handlers, m.newDelayEventHandler(m, topic, fn))
	}
	return nil
}

func (m *delayEventModuleImpl) GetHandlers() []delayEventHandler {
	return m.handlers
}

func (m *delayEventModuleImpl) GetAckTimeout() time.Duration {
	return m.ackTimeout
}

func (m *delayEventModuleImpl) GetBackOffPolicy() *backoff.ExponentialBackOff {
	return m.backoffPolicy
}

func (m *delayEventModuleImpl) newDelayEventHandler(module DelayEventModule, topic string, fn DelayEventFunction) delayEventHandler {
	return &delayEventHandlerImpl{
		module: module,
		topic:  topic,
		fn:     fn,
	}
}

// NewExponentialBackOff creates an instance of ExponentialBackOff using default values.
func getDefaultBackOffPolicy() *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     2 * time.Second,
		RandomizationFactor: 0.5,
		Multiplier:          2,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      30 * time.Second,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}
