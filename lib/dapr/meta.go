package dapr

import (
	"context"
	"github.com/dapr/dapr/pkg/api/grpc/metadata"
	"github.com/spf13/cast"
)

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

type MetaManager interface {
	GetAppId() string
	GetApiVersion() int
	GetMetaValue(key string) string
	GetMetaValues(key string) []string
}

type metaManagerImpl struct {
	ctx context.Context
}

func NewMetaManager(ctx context.Context) MetaManager {
	return &metaManagerImpl{ctx: ctx}
}

func (m metaManagerImpl) GetAppId() string {
	return m.GetMetaValue(MetaAppId)
}

func (m metaManagerImpl) GetApiVersion() int {
	return cast.ToInt(m.GetMetaValue(MetaApiVersion))
}

// GetMetaValues get grpc meta values
func (m metaManagerImpl) GetMetaValues(key string) []string {
	md, ok := metadata.FromIncomingContext(m.ctx)
	if !ok {
		return nil
	}
	return md.Get(key)
}

// GetMetaValue get the first grpc meta value
func (m metaManagerImpl) GetMetaValue(key string) string {
	md, ok := metadata.FromIncomingContext(m.ctx)
	if !ok {
		return ""
	}
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
