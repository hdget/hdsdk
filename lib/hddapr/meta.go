package hddapr

const (
	MetaPrefix  = "Hd-"
	MetaAppId   = MetaPrefix + "Application"
	MetaVersion = MetaPrefix + "Version"
)

var (
	// AllMetaKeys 所有meta的关键字
	AllMetaKeys = []string{
		MetaAppId,
		MetaVersion,
	}
)
