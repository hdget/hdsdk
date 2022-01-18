package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/hdget/hdsdk/lib/wx/wxmp/cache"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

type Wxmp interface {
	// 校验凭证
	Auth(appId, appSecret, code string) (*typwx.WxmpSession, error)
	Decrypt(appId, sessionKey, encryptedData, iv string) (*typwx.WxmpUserInfo, error)
}

type implWxmp struct{}

var (
	_      Wxmp = (*implWxmp)(nil)
	_cache cache.Cache
)

func init() {
	_cache = cache.New()
}

func New() Wxmp {
	return &implWxmp{}
}

// Auth 登录凭证校验，获取openid,unionid和session_key
func (w *implWxmp) Auth(appId, appSecret, code string) (*typwx.WxmpSession, error) {
	// url to get sessionKey, openId and unionId from Weixin server
	// do http get request against Wechat server
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appId,
		appSecret,
		code,
	)

	// 登录凭证校验
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "auth with wechat server")
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error auth with wechat server, status_code: %d", resp.StatusCode())
	}

	var session typwx.WxmpSession
	err = json.Unmarshal(resp.Body(), &session)
	if err != nil {
		return nil, err
	}

	// 保存到缓存中
	err = _cache.SetSessKey(appId, session.UnionId, session.SessionKey)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (w *implWxmp) Decrypt(appId, wxId, encryptedData, iv string) (*typwx.WxmpUserInfo, error) {
	sessKey, err := _cache.GetSessKey(appId, wxId)
	if err != nil {
		return nil, errors.Wrap(err, "session key not found, you should invoke wx.login() firstly")
	}

	return decrypt(appId, sessKey, encryptedData, iv)
}
