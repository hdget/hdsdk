package dapr

import (
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	"go/importer"
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
)

func NewGrpcServer(logger intf.LoggerProvider, address string, modulePaths ...string) (Server, error) {
	service, err := grpc.NewService(address)
	if err != nil {
		return nil, errors.Wrap(err, "new dapr grpc server")
	}

	srv := &serverImpl{Service: service, logger: logger}
	if err = srv.initialize(modulePaths...); err != nil {
		return nil, err
	}

	return srv, nil
}

func NewHttpServer(logger intf.LoggerProvider, address string, modulePaths ...string) (Server, error) {
	service, err := grpc.NewService(address)
	if err != nil {
		return nil, errors.Wrap(err, "new dapr http server")
	}

	srv := &serverImpl{Service: service, logger: logger}
	if err = srv.initialize(modulePaths...); err != nil {
		return nil, err
	}
	return srv, nil
}

// LoadInvocationModules 获取所有服务调用模块, args为服务模块所在的文件路径
func LoadInvocationModules(invocationModulePath string) map[string]InvocationModule {
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

// GetBindingHandlers todo:需要通过反射获取bindingHandlers
func (impl *serverImpl) GetBindingHandlers() map[string]common.BindingInvocationHandler {
	return nil
}

func registerInvocationModule(module InvocationModule) {
	_moduleName2invocationModule[module.GetMeta().ModuleName] = module
}

func registerEventModule(module EventModule) {
	_moduleName2eventModule[module.GetMeta().ModuleName] = module
}
