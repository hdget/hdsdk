package svc

import (
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/lib/dapr"
	"github.com/pkg/errors"
)

type DaprService interface {
	Service
	GetInvocationHandlers() map[string]common.ServiceInvocationHandler
	GetBindingHandlers() map[string]common.BindingInvocationHandler
	GetEvents() []dapr.Event
}

type baseDaprService struct {
}

func NewDaprService() DaprService {
	return &baseDaprService{}
}

func (impl *baseDaprService) GetInvocationHandlers() map[string]common.ServiceInvocationHandler {
	// 获取handlers
	handlers := make(map[string]common.ServiceInvocationHandler)
	for _, module := range GetInvocationModules() {
		for name, anyHandler := range module.GetHandlers() {
			// assert common.ServiceInvocationHandler
			handler, ok := anyHandler.(common.ServiceInvocationHandler)
			if ok {
				handlers[name] = handler
			}
		}
	}

	return handlers
}

// Initialize 初始化server
func (impl *baseDaprService) Initialize(server any, generators ...Generator) error {
	daprServer, ok := server.(common.Service)
	if !ok {
		return errors.New("invalid dapr common.service")
	}

	for method, handler := range impl.GetInvocationHandlers() {
		if err := daprServer.AddServiceInvocationHandler(method, handler); err != nil {
			return errors.Wrap(err, "adding invocation handler")
		}
	}

	for name, handler := range impl.GetBindingHandlers() {
		if err := daprServer.AddBindingInvocationHandler(name, handler); err != nil {
			return errors.Wrap(err, "adding binding handler")
		}
	}

	for _, event := range impl.GetEvents() {
		if err := daprServer.AddTopicEventHandler(event.Subscription, event.Handler); err != nil {
			return errors.Wrap(err, "adding event handler")
		}
	}

	// 注册生成的依赖文件
	for _, gen := range generators {
		err := gen.Register()
		if err != nil {
			return err
		}
	}

	return nil
}

func (impl *baseDaprService) GetBindingHandlers() map[string]common.BindingInvocationHandler {
	return nil
}

func (impl *baseDaprService) GetEvents() []dapr.Event {
	return nil
}
