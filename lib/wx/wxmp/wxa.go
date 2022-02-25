package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/pkg/errors"
)

type Param func(param *typwx.CommonWxaCodeParam)

// CreateLimitedWxaCode 创建小程序码
func (impl *implWxmp) CreateLimitedWxaCode(appId, appSecret, path string, args ...Param) (interface{}, error) {
	accessToken, err := impl.getAccessToken(appId, appSecret)
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &typwx.LimitedWxaCodeParam{
		Path: path,
		CommonWxaCodeParam: &typwx.CommonWxaCodeParam{
			EnvVersion: "release",
		},
	}
	for _, arg := range args {
		arg(body.CommonWxaCodeParam)
	}

	client := resty.New()
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacode?access_token=%s", accessToken)
	resp, err := client.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	var result typwx.WxmpWxaCodeResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.New("invalid wxmp wxa code result")
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	return result.Buffer, nil
}

// CreateUnLimitedWxaCode 创建小程序码
func (impl *implWxmp) CreateUnLimitedWxaCode(appId, appSecret, scene, page string, args ...Param) (interface{}, error) {
	accessToken, err := impl.getAccessToken(appId, appSecret)
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &typwx.UnLimitedWxaCodeParam{
		Scene: scene,
		Page:  page,
		CommonWxaCodeParam: &typwx.CommonWxaCodeParam{
			EnvVersion: "release",
		},
	}
	for _, arg := range args {
		arg(body.CommonWxaCodeParam)
	}

	client := resty.New()
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", accessToken)
	resp, err := client.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	var result typwx.WxmpWxaCodeResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.New("invalid wxmp wxa code result")
	}

	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}

	return result.Buffer, nil
}

func Trial(param *typwx.CommonWxaCodeParam) {
	param.EnvVersion = "trial"
}

func Develop(param *typwx.CommonWxaCodeParam) {
	param.EnvVersion = "develop"
}
