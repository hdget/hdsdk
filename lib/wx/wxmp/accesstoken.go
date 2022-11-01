package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/pkg/errors"
)

func (impl *implWxmp) getAccessToken(appId, appSecret string) (string, error) {
	// 尝试从缓存中获取access token
	cachedAccessToken, _ := _cache.GetAccessToken(appId)
	if cachedAccessToken != "" {
		return cachedAccessToken, nil
	}

	// 如果从缓存中获取不到，尝试请求access token
	wxAccessToken, err := impl.requestAccessToken(appId, appSecret)
	if err != nil {
		return "", err
	}

	err = _cache.SetAccessToken(appId, wxAccessToken.AccessToken, wxAccessToken.ExpiresIn-1000)
	if err != nil {
		return "", err
	}

	return wxAccessToken.AccessToken, nil
}

func (impl *implWxmp) requestAccessToken(appId, appSecret string) (*typwx.WxmpAccessToken, error) {
	wxAccessTokenTmpl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	wxAccessTokenURL := fmt.Sprintf(wxAccessTokenTmpl, appId, appSecret)

	client := resty.New()
	resp, err := client.R().Get(wxAccessTokenURL)
	if err != nil {
		return nil, err
	}

	var result typwx.WxmpAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.New("invalid wxmp access token result")
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	if result.AccessToken == "" {
		return nil, fmt.Errorf("empty access token, url: %s, resp: %s", wxAccessTokenURL, string(resp.Body()))
	}

	return &result.WxmpAccessToken, nil
}
