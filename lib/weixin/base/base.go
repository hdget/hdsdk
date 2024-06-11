package base

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/v2"
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
	urlGetUnionId       = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
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

// GetUser 通过openId获取用户信息
func (b *ApiWeixin) GetUser(openid string) (*types.UserInfo, error) {
	accessToken, err := b.GetAccessToken()
	if err != nil {
		hdsdk.Logger().Error("get access token", "err", err)
		return nil, err
	}

	url := fmt.Sprintf(urlGetUnionId, accessToken, openid)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "get access token, appId: %s", b.AppId)
	}

	var result types.UserInfo
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal result, body: %s", convert.BytesToString(resp.Body()))
	}

	return &result, nil
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

	resp, err := resty.New().R().Get(wxAccessTokenURL)
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
