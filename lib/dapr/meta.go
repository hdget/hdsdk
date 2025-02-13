package dapr

import (
	"context"
	"github.com/hdget/hdsdk/v2/lib/code"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

const (
	MetaKeyAppId   = "Hd-App-Id"
	MetaKeyRelease = "Hd-Release"
	MetaKeyTid     = "Hd-Tid"
	MetaKeyEtid    = "Hd-Etid"
	MetaKeyEuid    = "Hd-Euid"  // encoded user id
	MetaKeyErids   = "Hd-Erids" // encoded role ids
	MetaKeyCaller  = "dapr-caller-app-id"
)

var (
	// MetaKeys 所有meta的关键字
	_httpHeaderKeys = []string{
		MetaKeyEtid,
		MetaKeyAppId,
		MetaKeyRelease,
	}
)

type MetaManager interface {
	GetHttpHeaderKeys() []string
	GetValue(ctx context.Context, key string) string
	GetValues(ctx context.Context, key string) []string
	GetAppId(ctx context.Context) string
	GetRelease(ctx context.Context) string
	GetCaller(ctx context.Context) string
	GetUserId(ctx context.Context) int64
	GetRoleIds(ctx context.Context) []int64
	GetTenantId(ctx context.Context) int64
	GetEtid(ctx context.Context) string
	// DEPRECATED
	OldGetRoles(ctx context.Context) []*Role
	OldGetRoleValues(ctx context.Context) []string
	OldGetRoleIds(ctx context.Context) []int64
	OldGetPermIds(ctx context.Context) []int64
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

func (m metaManagerImpl) GetRelease(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyRelease)
}

func (m metaManagerImpl) GetCaller(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyCaller)
}

func (m metaManagerImpl) GetRoleIds(ctx context.Context) []int64 {
	return code.New().DecodeInt64Slice(m.GetValue(ctx, MetaKeyErids))
}

func (m metaManagerImpl) GetUserId(ctx context.Context) int64 {
	return code.New().DecodeInt64(m.GetValue(ctx, MetaKeyEuid))
}

func (m metaManagerImpl) GetTenantId(ctx context.Context) int64 {
	if v := m.GetValue(ctx, MetaKeyTid); v != "" {
		return cast.ToInt64(v)
	}
	return code.New().DecodeInt64(m.GetValue(ctx, MetaKeyEtid))
}

func (m metaManagerImpl) GetEtid(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyEtid)
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
