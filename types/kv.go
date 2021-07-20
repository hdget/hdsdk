package types

type KvProvider interface {
}

// key/value storage ability
const (
	_                    = SdkCategoryKv + iota
	CAP_PROVIDER_KV_ETCD // etcd能力
)
