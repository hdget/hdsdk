package hdsdk

type Option func(*SdkInstance)

// WithConfigVar if config var set then load config to config var
func WithConfigVar(configVar any) Option {
	return func(i *SdkInstance) {
		i.configVar = configVar
	}
}

func WithDebug() Option {
	return func(i *SdkInstance) {
		i.debug = true
	}
}
