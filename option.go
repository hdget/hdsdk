package hdsdk

type optionObject struct {
	configFilePath string
	debug          bool // debug mode
}

type Option func(*optionObject)

var (
	defaultSdkOption = &optionObject{
		debug: false,
	}
)

func WithDebug() Option {
	return func(o *optionObject) {
		o.debug = true
	}
}

func WithConfigFile(configFilePath string) Option {
	return func(o *optionObject) {
		o.configFilePath = configFilePath
	}
}
