package cache

import (
	"fmt"
	"github.com/hdget/hdsdk"
)

const (
	WXMP_SESSION_KEY         = `wxmp:%s:session`
	WXMP_SESSION_KEY_EXPIRES = 60 // session key过期时间60秒
)

type Cache interface {
	GetSessKey(appId string) (string, error)
	SetSessKey(appId, sessKey string) error
}

type implCache struct{}

var _ Cache = (*implCache)(nil)

func New() Cache {
	return &implCache{}
}

func (impl *implCache) GetSessKey(appId string) (string, error) {
	return hdsdk.Redis.My().GetString(getSessionKey(appId))
}

func (impl *implCache) SetSessKey(appId, sessKey string) error {
	return hdsdk.Redis.My().SetEx(getSessionKey(appId), sessKey, WXMP_SESSION_KEY_EXPIRES)
}

func getSessionKey(appId string) string {
	return fmt.Sprintf(WXMP_SESSION_KEY, appId)
}
