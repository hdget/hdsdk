package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/v2/lib/weixin/base"
	"github.com/hdget/hdsdk/v2/lib/weixin/types"
	"github.com/pkg/errors"
)

type ApiWxmp interface {
	GetWxId(code string) (string, string, error) // 校验凭证
	DecryptUserInfo(encryptedData, iv string) (*UserInfo, error)
	DecryptMobileInfo(encryptedData, iv string) (*MobileInfo, error)
	CreateLimitedWxaCode(path string, width int, options ...WxaCodeOption) ([]byte, error)
	CreateUnLimitedWxaCode(scene, page string, width int, options ...WxaCodeOption) ([]byte, error)
	GetUserPhoneNumber(code string) (*MobileInfo, error)
	GetAccessToken() (string, error)
}

type wxmpImpl struct {
	*base.ApiWeixin
}

var (
	_ ApiWxmp = (*wxmpImpl)(nil)
)

const (
	wxmpSessionKeyExpires = 3600 // session key过期时间3600秒
	urlCode2Session       = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	urlGetUserPhoneNumber = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
)

func New(appId, appSecret string) ApiWxmp {
	return &wxmpImpl{
		ApiWeixin: base.New(types.WeixinAppWxmp, appId, appSecret),
	}
}

// GetWxId 登录凭证校验，获取openid,unionid
func (impl *wxmpImpl) GetWxId(code string) (string, string, error) {
	// url to get sessionKey, openId and unionId from Weixin server
	// do http get request against Wechat server
	url := fmt.Sprintf(urlCode2Session, impl.AppId, impl.AppSecret, code)

	// 登录凭证校验
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return "", "", errors.Wrap(err, "wxmp code to session")
	}
	if resp.StatusCode() != 200 {
		return "", "", fmt.Errorf("wxmp code to session, status_code: %d", resp.StatusCode())
	}

	var session Session
	err = impl.ParseResult(resp.Body(), &session)
	if err != nil {
		return "", "", errors.Wrap(err, "unmarshal weixin api response")
	}

	if session.SessionKey == "" {
		return "", "", errors.New("empty session key")
	}

	// 保存到缓存中
	err = impl.Cache.SetSessKey(session.SessionKey, wxmpSessionKeyExpires)
	if err != nil {
		return "", "", err
	}

	return session.OpenId, session.UnionId, nil
}

func (impl *wxmpImpl) GetUserPhoneNumber(code string) (*MobileInfo, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "get access token")
	}

	body := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	resp, err := resty.New().R().SetBody(body).Post(fmt.Sprintf(urlGetUserPhoneNumber, accessToken))
	if err != nil {
		return nil, err
	}

	var result GetUserPhoneNumberResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.New("invalid wxmp access token result")
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	return &result.PhoneInfo, nil
}
