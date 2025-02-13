package dapr

import (
	"context"
	"github.com/spf13/cast"
	"strings"
)

type Role struct {
	Name  string // 角色名
	Level int    // 角色级别
}

// DEPRECATED: compatible purpose
const (
	MetaKeyUserId     = "Hd-User-Id"
	MetaKeyRoleValues = "Hd-Role-Values"
	MetaKeyRoleIds    = "Hd-Role-Ids"
	MetaKeyPermIds    = "Hd-Perm-Ids"
)

func (m metaManagerImpl) OldGetRoles(ctx context.Context) []*Role {
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

func (m metaManagerImpl) OldGetRoleValues(ctx context.Context) []string {
	return m.GetValues(ctx, MetaKeyRoleValues)
}

func (m metaManagerImpl) OldGetUserId(ctx context.Context) int64 {
	return cast.ToInt64(m.GetValue(ctx, MetaKeyUserId))
}

func (m metaManagerImpl) OldGetRoleIds(ctx context.Context) []int64 {
	roleIds := make([]int64, 0)
	for _, v := range m.GetValues(ctx, MetaKeyRoleIds) {
		roleIds = append(roleIds, cast.ToInt64(v))
	}
	return roleIds
}

func (m metaManagerImpl) OldGetPermIds(ctx context.Context) []int64 {
	permIds := make([]int64, 0)
	for _, v := range m.GetValues(ctx, MetaKeyPermIds) {
		permIds = append(permIds, cast.ToInt64(v))
	}
	return permIds
}
