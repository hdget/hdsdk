package types

type NosqlProvider interface {
}

// nosql db ability
const (
	_                     = SdkCategoryNosql + iota
	CAP_PROVIDER_NOSQL_ES // elasticSearch能力
)
