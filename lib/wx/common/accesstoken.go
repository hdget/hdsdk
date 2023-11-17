package common

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/hdutils"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/pkg/errors"
)

const (
	tplUrlGetWxAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func RequestAccessToken(appId, appSecret string) (*typwx.WxAccessToken, error) {
	wxAccessTokenURL := fmt.Sprintf(tplUrlGetWxAccessToken, appId, appSecret)

	client := resty.New()
	resp, err := client.R().Get(wxAccessTokenURL)
	if err != nil {
		return nil, errors.Wrapf(err, "get access token, appId: %s", appId)
	}

	var result typwx.WxAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal result, body: %s", hdutils.BytesToString(resp.Body()))
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	if result.AccessToken == "" {
		return nil, fmt.Errorf("empty access token, url: %s, resp: %s", wxAccessTokenURL, hdutils.BytesToString(resp.Body()))
	}

	return &result.WxAccessToken, nil
}
