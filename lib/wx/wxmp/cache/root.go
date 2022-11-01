package cache

import (
	"fmt"
	"hdsdk"
)

const (
	WXMP_SESSION_KEY         = `wxmp:%s:session`
	WXMP_ACCESS_TOKEN        = `wxmp:%s:accesstoken`
	WXMP_SESSION_KEY_EXPIRES = 3600 // session key过期时间3600秒
)

type Cache interface {
	GetSessKey(appId string) (string, error)
	SetSessKey(appId, sessKey string) error
	GetAccessToken(appId string) (string, error)
	SetAccessToken(appId, accessToken string, expires int) error
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

func (impl *implCache) GetAccessToken(appId string) (string, error) {
	bs, err := hdsdk.Redis.My().Get(getAccessToken(appId))
	return string(bs), err
}

func (impl *implCache) SetAccessToken(appId, token string, expires int) error {
	return hdsdk.Redis.My().SetEx(getAccessToken(appId), token, expires)
}

func getAccessToken(appId string) string {
	return fmt.Sprintf(WXMP_ACCESS_TOKEN, appId)
}

func getSessionKey(appId string) string {
	return fmt.Sprintf(WXMP_SESSION_KEY, appId)
}
