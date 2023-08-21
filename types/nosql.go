package types

type NosqlProvider interface {
}

// nosql db ability
const (
	_                             = SdkCategoryNosql + iota
	CapProviderNosqlElasticSearch // elasticSearch能力
)
