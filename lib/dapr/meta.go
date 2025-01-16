package dapr

import (
	"context"
	"github.com/hdget/hdsdk/v2/lib/code"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

type Role struct {
	Name  string // 角色名
	Level int    // 角色级别
}

const (
	MetaTenantId   = "Hd-Tid"
	MetaKeyAppId   = "Hd-App-Id"
	MetaKeyRelease = "Hd-Release"
	MetaKeyEuid    = "Hd-Euid"  // encoded user id
	MetaKeyErids   = "Hd-Erids" // encoded role ids
	MetaKeyCaller  = "dapr-caller-app-id"
)

var (
	// MetaKeys 所有meta的关键字
	_httpHeaderKeys = []string{
		MetaTenantId,
		MetaKeyAppId,
		MetaKeyRelease,
	}
)

type MetaManager interface {
	GetHttpHeaderKeys() []string
	GetValue(ctx context.Context, key string) string
	GetValues(ctx context.Context, key string) []string
	GetTenantId(ctx context.Context) int64
	GetAppId(ctx context.Context) string
	GetRelease(ctx context.Context) string
	GetCaller(ctx context.Context) string
	GetUserId(ctx context.Context, secret []byte) int64
	GetRoleIds(ctx context.Context, secret []byte) []int64
}

type metaManagerImpl struct {
}

func Meta() MetaManager {
	return &metaManagerImpl{}
}

func (m metaManagerImpl) GetTenantId(ctx context.Context) int64 {
	return cast.ToInt64(m.GetValue(ctx, MetaTenantId))
}

func (m metaManagerImpl) GetAppId(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyAppId)
}

func (m metaManagerImpl) GetHttpHeaderKeys() []string {
	return _httpHeaderKeys
}

func (m metaManagerImpl) GetRelease(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyRelease)
}

func (m metaManagerImpl) GetCaller(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyCaller)
}

func (m metaManagerImpl) GetRoleIds(ctx context.Context, secret []byte) []int64 {
	return code.New().DecodeInt64Slice(m.GetValue(ctx, MetaKeyErids), secret)
}

func (m metaManagerImpl) GetUserId(ctx context.Context, secret []byte) int64 {
	return code.New().DecodeInt64(m.GetValue(ctx, MetaKeyEuid), secret)
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
