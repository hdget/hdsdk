package svc

type Service interface {
	Initialize(server any, generators ...Generator) error
}

var (
	registryInvocationModule = make(map[string]InvocationModule)
)

func GetInvocationModules() map[string]InvocationModule {
	return registryInvocationModule
}

func addInvocationModule(moduleName string, module InvocationModule) {
	registryInvocationModule[moduleName] = module
}
