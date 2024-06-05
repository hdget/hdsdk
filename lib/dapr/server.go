package dapr

import (
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/dapr/go-sdk/service/http"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	"go/importer"
	"net"
)

type Server interface {
	Start() error
	Stop() error
	GracefulStop() error
	GetInvocationHandlers() map[string]common.ServiceInvocationHandler
	GetBindingHandlers() map[string]common.BindingInvocationHandler
	GetEvents() []Event
}

type serverImpl struct {
	common.Service
	logger intf.LoggerProvider
}

var (
	_moduleName2invocationModule = make(map[string]InvocationModule) // service invocation module
	_moduleName2eventModule      = make(map[string]EventModule)      // topic event module
	_moduleName2healthModule     = make(map[string]HealthModule)     // health module
)

func NewGrpcServer(logger intf.LoggerProvider, address string, modulePaths ...string) (Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("grpc server failed to listen on %s: %w", address, err)
	}

	// install health check handler
	grpcServer := grpc.NewServiceWithListener(lis)

	srv := &serverImpl{Service: grpcServer, logger: logger}
	if err = srv.initialize(modulePaths...); err != nil {
		return nil, err
	}

	return srv, nil
}

func NewHttpServer(logger intf.LoggerProvider, address string, modulePaths ...string) (Server, error) {
	httpServer := http.NewServiceWithMux(address, nil)

	srv := &serverImpl{Service: httpServer, logger: logger}
	if err := srv.initialize(modulePaths...); err != nil {
		return nil, err
	}
	return srv, nil
}

// loadInvocationModules 获取所有服务调用模块, args为服务模块所在的文件路径
func loadInvocationModules(invocationModulePath string) map[string]InvocationModule {
	_, _ = importer.Default().Import(invocationModulePath)
	return _moduleName2invocationModule
}

func (impl *serverImpl) Start() error {
	return impl.Service.Start()
}

func (impl *serverImpl) Stop() error {
	return impl.Service.Stop()
}

func (impl *serverImpl) GracefulStop() error {
	return impl.Service.GracefulStop()
}

// Initialize 初始化server
func (impl *serverImpl) initialize(modulePaths ...string) error {
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
		if err := impl.AddTopicEventHandler(event.Subscription, event.Handler); err != nil {
			return errors.Wrap(err, "adding event handler")
		}
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

func (impl *serverImpl) GetEvents() []Event {
	// 获取handlers
	events := make([]Event, 0)
	for _, m := range _moduleName2eventModule {
		for _, h := range m.GetHandlers() {
			events = append(events, GetEvent(m.GetPubSub(), h.GetTopic(), h.GetEventFunction(impl.logger)))
		}
	}
	return events
}

func (impl *serverImpl) GetHealthCheckHandler() common.HealthCheckHandler {
	if len(_moduleName2healthModule) == 0 {
		return EmptyHealthCheckFunction
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
	_moduleName2invocationModule[module.GetInfo().ModuleName] = module
}

func registerEventModule(module EventModule) {
	_moduleName2eventModule[module.GetInfo().ModuleName] = module
}

func registerHealthModule(module HealthModule) {
	_moduleName2healthModule[module.GetInfo().ModuleName] = module
}
