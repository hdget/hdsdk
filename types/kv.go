package types

type KvProvider interface {
}

// key/value storage ability
const (
	_                 = SdkCategoryKv + iota
	CapProviderKvEtcd // etcd能力
)
