package svc

const (
	MetaPrefix  = "Hd-"
	MetaAppId   = MetaPrefix + "App-Id"
	MetaVersion = MetaPrefix + "Version"
)

var (
	// AllMetaKeys 所有meta的关键字
	AllMetaKeys = []string{
		MetaAppId,
		MetaVersion,
	}
)
