package hdsdk

type Option func(*SdkInstance)

func WithDebug() Option {
	return func(i *SdkInstance) {
		i.debug = true
	}
}
