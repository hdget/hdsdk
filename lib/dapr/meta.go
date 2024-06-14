package dapr

import (
	"context"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
	"strings"
)

type RoleValue struct {
	Id    int64
	Level int
}

const (
	MetaKeyAppId      = "Hd-App-Id"
	MetaKeyRelease    = "Hd-Release"
	MetaKeyUserId     = "Hd-User-Id"
	MetaKeyRoleValues = "Hd-Role-Values"
	MetaKeyPermIds    = "Hd-Perm-Ids"
)

var (
	// MetaKeys 所有meta的关键字
	_httpHeaderKeys = []string{
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
	GetUserId(ctx context.Context) int64
	GetRoleValues(ctx context.Context) []*RoleValue
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

func (m metaManagerImpl) GetRelease(ctx context.Context) string {
	return m.GetValue(ctx, MetaKeyRelease)
}

func (m metaManagerImpl) GetUserId(ctx context.Context) int64 {
	return cast.ToInt64(m.GetValue(ctx, MetaKeyUserId))
}

func (m metaManagerImpl) GetRoleValues(ctx context.Context) []*RoleValue {
	results := make([]*RoleValue, 0)
	for _, v := range m.GetValues(ctx, MetaKeyRoleValues) {
		tokens := strings.Split(v, ":")
		if len(tokens) != 2 {
			return nil
		}

		roleId := cast.ToInt64(tokens[0])
		if roleId == 0 {
			return nil
		}

		results = append(results, &RoleValue{
			Id:    roleId,
			Level: cast.ToInt(tokens[1]),
		})
	}
	return results
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
