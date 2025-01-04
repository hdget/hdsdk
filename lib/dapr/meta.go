package dapr

import (
	"context"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
	"strings"
)

type Role struct {
	Name  string // 角色名
	Level int    // 角色级别
}

const (
	MetaTenantId      = "Hd-Tid"
	MetaKeyAppId      = "Hd-App-Id"
	MetaKeyRelease    = "Hd-Release"
	MetaKeyUserId     = "Hd-User-Id"
	MetaKeyRoleValues = "Hd-Role-Values"
	MetaKeyPermIds    = "Hd-Perm-Ids"
	MetaKeyRoleIds    = "Hd-Role-Ids"
	MetaKeyCaller     = "dapr-caller-app-id"
	MetaKeyUid        = "Hd-Uid"
	MetaKeyRid        = "Hd-Rid"
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
	GetUserId(ctx context.Context) int64
	GetRoles(ctx context.Context) []*Role
	GetRoleValues(ctx context.Context) []string
	GetPermIds(ctx context.Context) []int64
	GetCaller(ctx context.Context) string
	GetUid(ctx context.Context) uint64
	GetRoleIds(ctx context.Context) []uint64
}

type metaManagerImpl struct {
	secret []byte
}

func Meta(secret ...byte) MetaManager {
	return &metaManagerImpl{
		secret: secret,
	}
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

func (m metaManagerImpl) GetUserId(ctx context.Context) int64 {
	return cast.ToInt64(m.GetValue(ctx, MetaKeyUserId))
}

func (m metaManagerImpl) GetCaller(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyCaller)
}

func (m metaManagerImpl) GetRoles(ctx context.Context) []*Role {
	roles := make([]*Role, 0)
	for _, roleValue := range m.GetValues(ctx, MetaKeyRoleValues) {
		tokens := strings.Split(roleValue, ":")
		if len(tokens) != 2 {
			return nil
		}
		roles = append(roles, &Role{
			Name:  tokens[0],
			Level: cast.ToInt(tokens[1]),
		})
	}
	return roles
}

func (m metaManagerImpl) GetRoleIds(ctx context.Context) []uint64 {
	roleIds, _ := Coder(m.secret).DecodeUint64Slice(m.GetValue(ctx, MetaKeyRid))
	return roleIds
}

func (m metaManagerImpl) GetUid(ctx context.Context) uint64 {
	uid, _ := Coder(m.secret).DecodeUint64(m.GetValue(ctx, MetaKeyUid))
	return uid
}

func (m metaManagerImpl) GetRoleValues(ctx context.Context) []string {
	return m.GetValues(ctx, MetaKeyRoleValues)
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
