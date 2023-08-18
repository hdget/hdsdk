package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/wx/common"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/hdget/hdsdk/lib/wx/wxmp/cache"
	"github.com/pkg/errors"
)

type Wxmp interface {
	Auth(appId, appSecret, code string) (*typwx.WxmpSession, error) // 校验凭证
	DecryptUserInfo(appId, encryptedData, iv string) (*typwx.WxmpUserInfo, error)
	DecryptMobileInfo(appId, encryptedData, iv string) (*typwx.WxmpMobileInfo, error)
	CreateLimitedWxaCode(appId, appSecret, path string, width int, args ...Param) ([]byte, error)
	CreateUnLimitedWxaCode(appId, appSecret, scene, page string, width int, args ...Param) ([]byte, error)
	GetUserPhoneNumber(appId, appSecret, code string) (*typwx.WxmpMobileInfo, error)
	GetAccessToken(appId, appSecret string) (string, error)
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

	session := &typwx.WxmpSession{}
	err = json.Unmarshal(resp.Body(), session)
	if session.SessionKey == "" {
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal to WxmpSession")
		}

		// 如果unmarshal请求消息错误,尝试获取错误信息
		var errResp typwx.WxErrResponse
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			hdsdk.Logger.Error("unmarshal wx err response", "err", err)
		}
		return nil, errors.New(errResp.ErrMsg)
	}

	// 保存到缓存中
	err = _cache.SetSessKey(appId, session.SessionKey)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (w *implWxmp) GetUserPhoneNumber(appId, appSecret, code string) (*typwx.WxmpMobileInfo, error) {
	accessToken, err := w.GetAccessToken(appId, appSecret)
	if err != nil {
		return nil, errors.Wrap(err, "get access token")
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken)
	client := resty.New()

	body := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	resp, err := client.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	var result typwx.WxmpMobileResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.New("invalid wxmp access token result")
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	return &result.PhoneInfo, nil
}

func (w *implWxmp) GetAccessToken(appId, appSecret string) (string, error) {
	// 尝试从缓存中获取access token
	cachedAccessToken, _ := _cache.GetAccessToken(appId)
	if cachedAccessToken != "" {
		return cachedAccessToken, nil
	}

	// 如果从缓存中获取不到，尝试请求access token
	wxAccessToken, err := common.RequestAccessToken(appId, appSecret)
	if err != nil {
		return "", err
	}

	err = _cache.SetAccessToken(appId, wxAccessToken.AccessToken, wxAccessToken.ExpiresIn-1000)
	if err != nil {
		return "", err
	}

	return wxAccessToken.AccessToken, nil
}
