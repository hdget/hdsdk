package dapr

const (
	MetaPrefix     = "hd-"
	MetaAppId      = MetaPrefix + "app-id"
	MetaApiVersion = MetaPrefix + "api-version"
)

var (
	// AllMetaKeys 所有meta的关键字
	AllMetaKeys = []string{
		MetaAppId,
		MetaApiVersion,
	}
)
