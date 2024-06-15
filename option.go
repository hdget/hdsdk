package hdsdk

type sdkOption struct {
	app            string
	env            string
	configFilePath string
	debug          bool // debug mode
}

type Option func(*sdkOption)

var (
	defaultSdkOption = &sdkOption{
		app:   "",
		env:   "test",
		debug: false,
	}
)

func WithDebug() Option {
	return func(o *sdkOption) {
		o.debug = true
	}
}

func WithConfigFile(configFilePath string) Option {
	return func(o *sdkOption) {
		o.configFilePath = configFilePath
	}
}
