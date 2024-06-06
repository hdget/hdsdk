package dapr

import (
	"context"
	"github.com/dapr/dapr/pkg/api/grpc/metadata"
	"github.com/spf13/cast"
)

const (
	MetaKeyAppId      = "appId"
	MetaKeyApiVersion = "apiVersion"
	MetaKeyUserId     = "userId"
	MetaKeyRoleIds    = "roleIds"
	MetaKeyPermIds    = "permIds"
)

var (
	// MetaKeys 所有meta的关键字
	_httpHeaderKeys = []string{
		MetaKeyAppId,
		MetaKeyApiVersion,
	}
)

type MetaManager interface {
	GetHttpHeaderKeys() []string
	GetValue(ctx context.Context, key string) string
	GetValues(ctx context.Context, key string) []string
	GetAppId(ctx context.Context) string
	GetApiVersion(ctx context.Context) string
	GetUserId(ctx context.Context) int64
	GetRoleIds(ctx context.Context) []int64
	GetPermIds(ctx context.Context) []int64
}

type metaManagerImpl struct {
}

func Meta() MetaManager {
	return &metaManagerImpl{}
}

func (m metaManagerImpl) GetAppId(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyAppId)
}

func (m metaManagerImpl) GetHttpHeaderKeys() []string {
	return _httpHeaderKeys
}

func (m metaManagerImpl) GetApiVersion(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyApiVersion)
}

func (m metaManagerImpl) GetUserId(ctx context.Context) int64 {
	return cast.ToInt64(m.GetValue(ctx, MetaKeyUserId))
}

func (m metaManagerImpl) GetRoleIds(ctx context.Context) []int64 {
	roleIds := make([]int64, 0)
	for _, v := range m.GetValues(ctx, MetaKeyRoleIds) {
		roleIds = append(roleIds, cast.ToInt64(v))
	}
	return roleIds
}

func (m metaManagerImpl) GetPermIds(ctx context.Context) []int64 {
	permIds := make([]int64, 0)
	for _, v := range m.GetValues(ctx, MetaKeyPermIds) {
		permIds = append(permIds, cast.ToInt64(v))
	}
	return permIds
}

// GetValues get grpc meta values
func (m metaManagerImpl) GetValues(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md.Get(key)
}

// GetValue get the first grpc meta value
func (m metaManagerImpl) GetValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
