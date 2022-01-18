package cache

import (
	"fmt"
	"github.com/hdget/hdsdk"
)

const (
	WXMP_SESSION_KEY = `wxmp:%s:session:%s`
	// session key过期时间60秒
	WXMP_SESSION_KEY_EXPIRES = 60
)

type Cache interface {
	GetSessKey(appId, wxId string) (string, error)
	SetSessKey(appId, wxId, sessKey string) error
}

type implCache struct{}

var _ Cache = (*implCache)(nil)

func New() Cache {
	return &implCache{}
}

func (impl *implCache) GetSessKey(appId, wxId string) (string, error) {
	return hdsdk.Redis.My().GetString(getSessionKey(appId, wxId))
}

func (impl *implCache) SetSessKey(appId, wxId, sessKey string) error {
	return hdsdk.Redis.My().SetEx(getSessionKey(appId, wxId), sessKey, WXMP_SESSION_KEY_EXPIRES)
}

func getSessionKey(appId, wxId string) string {
	return fmt.Sprintf(WXMP_SESSION_KEY, appId, wxId)
}
