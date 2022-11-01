package wxoa

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"hdsdk"
	"hdsdk/lib/wx/typwx"
)

const (
	CACHE_KEY_WS_JSSDK_ACCESSTOKEN = "wxjssdk:access_token"
)

// WxoaAccessToken 类型
type WxoaAccessToken struct {
	Value     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

// nolint:errcheck
func (w *implWxoa) getAccessToken(appId, appSecret string) (string, error) {
	cachedAccessToken, err := _cache.GetAccessToken(appId)
	if err != nil {
		return "", errors.Wrap(err, "get wxoa cached access token")
	}

	if cachedAccessToken != "" {
		return cachedAccessToken, nil
	}

	accessToken, err := w.requestAccessToken(appId, appSecret)
	if err != nil {
		return "", err
	}

	err = _cache.SetAccessToken(appId, accessToken.Value, accessToken.ExpiresIn-1000)
	if err != nil {
		return "", errors.Wrap(err, "set wxoa access token to cache")
	}

	return accessToken.Value, nil
}

// 从微信服务器获取access token
func (w *implWxoa) requestAccessToken(appId, appSecret string) (*WxoaAccessToken, error) {
	wxUserAccessTokenTmpl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	wxUserAccessTokenURL := fmt.Sprintf(wxUserAccessTokenTmpl, appId, appSecret)

	client := resty.New()
	resp, err := client.R().Get(wxUserAccessTokenURL)
	if err != nil {
		return nil, err
	}

	token := &WxoaAccessToken{}
	err = json.Unmarshal(resp.Body(), token)
	if token.Value == "" {
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal to WxoaAccessToken")
		}

		// 如果unmarshal请求消息错误,尝试获取错误信息
		var errResp typwx.WxErrResponse
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			hdsdk.Logger.Error("unmarshal wx err response", "err", err)
		}
		return nil, fmt.Errorf("empty access token, url: %s, resp: %s", wxUserAccessTokenURL, string(resp.Body()))
	}

	return token, nil
}
