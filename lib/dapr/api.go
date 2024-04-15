package dapr

type APIer interface {
	Invoke(appId string, moduleVersion int, module, method string, data any, args ...string) ([]byte, error)
	GetServiceInvocationName(moduleVersion int, moduleName, handler string) string
}

type apiImpl struct {
}

func Api() APIer {
	return &apiImpl{}
}
