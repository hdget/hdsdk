package dapr

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/go-sdk/service/http"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	"net"
)

type Server interface {
	Start() error
	Stop() error
	GracefulStop() error
	GetInvocationHandlers() map[string]common.ServiceInvocationHandler
	GetBindingHandlers() map[string]common.BindingInvocationHandler
	GetEvents() []daprEvent
}

type serverImpl struct {
	common.Service
	logger intf.LoggerProvider
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	_moduleName2invocationModule = make(map[string]InvocationModule) // service invocation module
	_moduleName2eventModule      = make(map[string]EventModule)      // topic event module
	_moduleName2healthModule     = make(map[string]HealthModule)     // health module
	_moduleName2delayEventModule = make(map[string]DelayEventModule) // delay event module
)

func NewGrpcServer(logger intf.LoggerProvider, address string) (Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("grpc server failed to listen on %s: %w", address, err)
	}

	// install health check handler
	grpcServer := grpc.NewServiceWithListener(lis)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &serverImpl{
		Service: grpcServer,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}

	if err = appServer.initialize(); err != nil {
		return nil, err
	}

	return appServer, nil
}

func NewHttpServer(logger intf.LoggerProvider, address string) (Server, error) {
	httpServer := http.NewServiceWithMux(address, nil)

	ctx, cancel := context.WithCancel(context.Background())
	appServer := &serverImpl{
		Service: httpServer,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}

	if err := appServer.initialize(); err != nil {
		return nil, err
	}

	return appServer, nil
}

func (impl *serverImpl) Start() error {
	return impl.Service.Start()
}

func (impl *serverImpl) Stop() error {
	impl.cancel()
	return impl.Service.Stop()
}

func (impl *serverImpl) GracefulStop() error {
	impl.cancel()
	return impl.Service.GracefulStop()
}

// Initialize 初始化server
func (impl *serverImpl) initialize() error {
	// 注册health check handler
	if err := impl.AddHealthCheckHandler("", impl.GetHealthCheckHandler()); err != nil {
		return errors.Wrap(err, "adding health check handler")
	}

	// 注册各种类型的handlers
	for method, handler := range impl.GetInvocationHandlers() {
		if err := impl.AddServiceInvocationHandler(method, handler); err != nil {
			return errors.Wrap(err, "adding invocation handler")
		}
	}

	for name, handler := range impl.GetBindingHandlers() {
		if err := impl.AddBindingInvocationHandler(name, handler); err != nil {
			return errors.Wrap(err, "adding binding handler")
		}
	}

	for _, event := range impl.GetEvents() {
		if err := impl.AddTopicEventHandler(event.subscription, event.handler); err != nil {
			return errors.Wrap(err, "adding event handler")
		}
	}

	err := impl.SubscribeDelayEvents()
	if err != nil {
		return errors.Wrap(err, "adding delay event handler")
	}

	return nil
}

func (impl *serverImpl) GetInvocationHandlers() map[string]common.ServiceInvocationHandler {
	// 获取handlers
	handlers := make(map[string]common.ServiceInvocationHandler)
	for _, invocationModule := range _moduleName2invocationModule {
		for _, h := range invocationModule.GetHandlers() {
			handlers[h.GetInvokeName()] = h.GetInvokeFunction(impl.logger)
		}
	}
	return handlers
}

func (impl *serverImpl) GetEvents() []daprEvent {
	// 获取handlers
	events := make([]daprEvent, 0)
	for _, m := range _moduleName2eventModule {
		for _, h := range m.GetHandlers() {
			events = append(events, getDaprEvent(m.GetPubSub(), h.GetTopic(), h.GetEventFunction(impl.logger)))
		}
	}
	return events
}

func (impl *serverImpl) SubscribeDelayEvents() error {
	topic2delayEventHandler := make(map[string]delayEventHandler)
	for _, m := range _moduleName2delayEventModule {
		for _, h := range m.GetHandlers() {
			topic2delayEventHandler[h.GetTopic()] = h
		}
	}

	if len(topic2delayEventHandler) == 0 {
		return nil
	}

	// if we configured delay event module, but no message queue configured raise error
	if hdsdk.Mq() == nil {
		return errors.New("sdk message queue not initialized")
	}

	subscriber, err := hdsdk.Mq().Subscriber()
	if err != nil {
		return errors.Wrap(err, "new message queue subscriber")
	}

	for _, h := range topic2delayEventHandler {
		go h.Handle(impl.ctx, impl.logger, subscriber)
	}
	return nil
}

func (impl *serverImpl) GetHealthCheckHandler() common.HealthCheckHandler {
	if len(_moduleName2healthModule) == 0 {
		return emptyHealthCheckFunction
	}

	var firstModule HealthModule
	for _, module := range _moduleName2healthModule {
		firstModule = module
		break
	}
	return firstModule.GetHandler()
}

// GetBindingHandlers todo:需要通过反射获取bindingHandlers
func (impl *serverImpl) GetBindingHandlers() map[string]common.BindingInvocationHandler {
	return nil
}

func registerInvocationModule(module InvocationModule) {
	_moduleName2invocationModule[module.GetModuleInfo().ModuleName] = module
}

func registerEventModule(module EventModule) {
	_moduleName2eventModule[module.GetModuleInfo().ModuleName] = module
}

func registerDelayEventModule(module DelayEventModule) {
	_moduleName2delayEventModule[module.GetModuleInfo().ModuleName] = module
}

func registerHealthModule(module HealthModule) {
	_moduleName2healthModule[module.GetModuleInfo().ModuleName] = module
}
