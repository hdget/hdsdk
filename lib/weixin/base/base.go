package base

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/lib/weixin/cache"
	"github.com/hdget/hdsdk/v2/lib/weixin/types"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
)

type ApiWeixin struct {
	Logger    intf.LoggerProvider
	App       types.WeixinApp
	AppId     string
	AppSecret string
	Cache     cache.ApiWeixinCache
}

// ErrResponse 微信的错误响应
type ErrResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

const (
	urlGetWxAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func New(app types.WeixinApp, appId, appSecret string) *ApiWeixin {
	return &ApiWeixin{
		App:       app,
		AppId:     appId,
		AppSecret: appSecret,
		Cache:     cache.New(app, appId),
	}
}

func (b *ApiWeixin) GetAccessToken() (string, error) {
	// 尝试从缓存中获取access token
	cachedAccessToken, _ := b.Cache.GetAccessToken()
	if cachedAccessToken != "" {
		return cachedAccessToken, nil
	}

	// 如果从缓存中获取不到，尝试请求access token
	wxAccessToken, err := b.generateAccessToken()
	if err != nil {
		return "", err
	}

	err = b.Cache.SetAccessToken(wxAccessToken.AccessToken, wxAccessToken.ExpiresIn-1000)
	if err != nil {
		return "", err
	}

	return wxAccessToken.AccessToken, nil
}

func (b *ApiWeixin) ParseResult(data []byte, result any) error {
	err := json.Unmarshal(data, result)
	if err != nil { // 如果出错，则尝试解析API错误
		// 如果unmarshal请求消息错误,尝试获取错误信息
		var errResp ErrResponse
		if err = json.Unmarshal(data, &errResp); err == nil {
			return errors.New(errResp.ErrMsg)
		}
		return err
	}
	return nil
}

func (b *ApiWeixin) generateAccessToken() (*types.WxAccessToken, error) {
	wxAccessTokenURL := fmt.Sprintf(urlGetWxAccessToken, b.AppId, b.AppSecret)

	client := resty.New()
	resp, err := client.R().Get(wxAccessTokenURL)
	if err != nil {
		return nil, errors.Wrapf(err, "get access token, appId: %s", b.AppId)
	}

	var result types.WxAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal result, body: %s", convert.BytesToString(resp.Body()))
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	if result.AccessToken == "" {
		return nil, fmt.Errorf("empty access token, url: %s, resp: %s", wxAccessTokenURL, convert.BytesToString(resp.Body()))
	}

	return &result.WxAccessToken, nil
}
