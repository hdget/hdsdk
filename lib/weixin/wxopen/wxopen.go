package wxopen

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/v2/lib/weixin/base"
	"github.com/hdget/hdsdk/v2/lib/weixin/types"
	"github.com/pkg/errors"
)

type ApiWxopen interface {
	GetWxId(code string) (string, string, error) // 微信开放平台扫码登录
}

type wxopenImpl struct {
	*base.ApiWeixin
}

var (
	_ ApiWxopen = (*wxopenImpl)(nil)
)

const (
	urlOAuth2AccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

func New(appId, appSecret string) ApiWxopen {
	return &wxopenImpl{
		ApiWeixin: base.New(types.WeixinAppWxqr, appId, appSecret),
	}
}

func (impl *wxopenImpl) GetWxId(code string) (string, string, error) {
	url := fmt.Sprintf(urlOAuth2AccessToken, impl.AppId, impl.AppSecret, code)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return "", "", err
	}

	var result WxqrResult
	err = impl.ParseResult(resp.Body(), &result)
	if err != nil {
		return "", "", errors.Wrap(err, "unmarshal weixin api response")
	}

	if result.AccessToken == "" {
		return "", "", fmt.Errorf("empty access token, url: %s, resp: %s", url, string(resp.Body()))
	}

	return result.OpenId, result.UnionId, nil
}
