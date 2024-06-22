package hotconfig

type Option func(*hotConfigManager)
type SaveFunction func(configName string, data []byte) (Transactor, error)
type LoadFunction func(configName string) ([]byte, error)

func WithSaveFunction(fn SaveFunction) Option {
	return func(hc *hotConfigManager) {
		hc.saveFunction = fn
	}
}

func WithLoadFunction(fn LoadFunction) Option {
	return func(hc *hotConfigManager) {
		hc.loadFunction = fn
	}
}

func WithConfigStore(configStore string) Option {
	return func(hc *hotConfigManager) {
		hc.daprConfigStore = configStore
	}
}
