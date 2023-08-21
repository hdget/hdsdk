package cache

import (
	"fmt"
	"github.com/hdget/hdsdk"
)

const (
	tplWxmpSessionKey     = `wxmp:%s:session`
	tplWxmpAccessToken    = `wxmp:%s:accesstoken`
	wxmpSessionKeyExpires = 3600 // session key过期时间3600秒
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
	return hdsdk.Redis.My().SetEx(getSessionKey(appId), sessKey, wxmpSessionKeyExpires)
}

func (impl *implCache) GetAccessToken(appId string) (string, error) {
	bs, err := hdsdk.Redis.My().Get(getAccessToken(appId))
	return string(bs), err
}

func (impl *implCache) SetAccessToken(appId, token string, expires int) error {
	return hdsdk.Redis.My().SetEx(getAccessToken(appId), token, expires)
}

func getAccessToken(appId string) string {
	return fmt.Sprintf(tplWxmpAccessToken, appId)
}

func getSessionKey(appId string) string {
	return fmt.Sprintf(tplWxmpSessionKey, appId)
}
