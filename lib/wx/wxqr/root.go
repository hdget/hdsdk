package wxqr

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk"
	"github.com/pkg/errors"
)

type Wxqr interface {
	GetWxId(appId, appSecret, code string) (string, string, error)
}

type implWxqr struct{}

// WxErrResponse 生成微信的AccessToken
type WxErrResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var (
	_ Wxqr = (*implWxqr)(nil)
)

func New() Wxqr {
	return &implWxqr{}
}

// WxqrResponse 类型
type WxqrResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
	Scope        string `json:"scope"`
}

func (w *implWxqr) GetWxId(appId, appSecret, code string) (string, string, error) {
	resp, err := w.getWxqrResponse(appId, appSecret, code)
	if err != nil {
		return "", "", err
	}

	return resp.OpenId, resp.UnionId, nil
}

func (w *implWxqr) getWxqrResponse(appId, appSecret, code string) (*WxqrResponse, error) {
	wxAccessTokenTmpl := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	wxAccessTokenURL := fmt.Sprintf(wxAccessTokenTmpl, appId, appSecret, code)

	client := resty.New()
	resp, err := client.R().Get(wxAccessTokenURL)
	if err != nil {
		return nil, err
	}

	var wxAccessToken WxqrResponse
	err = json.Unmarshal(resp.Body(), &wxAccessToken)
	if err != nil {
		var errResp WxErrResponse
		// 如果unmarshal请求消息错误,尝试获取错误信息
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			hdsdk.Logger.Error("unmarshal wx err response", "err", err)
			return nil, errors.New(errResp.ErrMsg)
		}
		return nil, err
	}

	if wxAccessToken.AccessToken == "" {
		return nil, fmt.Errorf("empty access token, url: %s, resp: %s", wxAccessTokenURL, string(resp.Body()))
	}

	return &wxAccessToken, nil
}

//
//func (w *implWxqr) getAccessToken(appId, appSecret, code string) (string, error) {
//	// 尝试从缓存中获取access token
//	cachedAccessToken, err := _cache.GetAccessToken()
//	if err != nil {
//		sdk.Log.Error("msg", "get cached access token", "err", err)
//	}
//	if cachedAccessToken != "" {
//		return cachedAccessToken, nil
//	}
//
//	// 如果从缓存中获取不到，尝试通过refresh token刷新access token
//	wxAccessToken, err := w.refreshAccessToken(appId)
//	if err != nil {
//		sdk.Log.Error("msg", "refresh token", "err", err)
//	}
//
//	if wxAccessToken != nil && wxAccessToken.AccessToken != "" {
//		err = _cache.SetAccessToken(wxAccessToken.AccessToken, wxAccessToken.ExpiresIn)
//		if err != nil {
//			sdk.Log.Error("msg", "save access token to cache", "err", err)
//		}
//		return wxAccessToken.AccessToken, nil
//	}
//
//	// 如果刷新也获取不到，尝试重新获取
//	wxAccessToken, err = w.getWxqrResponse(appId, appSecret, code)
//	if err != nil {
//		return "", err
//	}
//
//	if wxAccessToken == nil || wxAccessToken.AccessToken == "" {
//		return "", errors.New("empty access token")
//	}
//
//	err = _cache.SetAccessToken(wxAccessToken.AccessToken, wxAccessToken.ExpiresIn)
//	if err != nil {
//		sdk.Log.Error("msg", "save access token to cache", "err", err)
//	}
//
//	err = _cache.SetRefreshToken(wxAccessToken.RefreshToken)
//	if err != nil {
//		sdk.Log.Error("msg", "save refresh token to cache", "err", err)
//	}
//	return wxAccessToken.AccessToken, nil
//}
//

//
//func (w *implWxqr) refreshAccessToken(appId string) (*WxqrResponse, error) {
//	refreshToken, err := _cache.GetRefreshToken()
//	if err != nil {
//		return nil, err
//	}
//
//	if refreshToken == "" {
//		return nil, errors.New("empty refresh token")
//	}
//
//	// 尝试请求获取新的access token
//	wxRefreshTokenTmpl := "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
//	wxRefreshTokenURL := fmt.Sprintf(wxRefreshTokenTmpl, appId, refreshToken)
//
//	client := resty.New()
//	resp, err := client.R().Get(wxRefreshTokenURL)
//	if err != nil {
//		return nil, err
//	}
//
//	var wxAccessToken WxqrResponse
//	err = json.Unmarshal(resp.Body(), &wxAccessToken)
//	if err != nil {
//		var errResp WxErrResponse
//		// 如果unmarshal请求消息错误,尝试获取错误信息
//		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
//			sdk.Log.Error("msg", "unmarshal wx err response", "err", err)
//			return nil, errors.New(errResp.ErrMsg)
//		}
//		return nil, err
//	}
//
//	if wxAccessToken.AccessToken == "" {
//		sdk.Log.Error("msg", "requested empty access token", "url", wxRefreshTokenURL, "resp", string(resp.Body()))
//		return nil, errors.New("empty access token")
//	}
//
//	return &wxAccessToken, nil
//}
