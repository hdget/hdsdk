package wxmp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk/hdutils"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/pkg/errors"
)

type Param func(param *typwx.CommonWxaCodeParam)

// CreateLimitedWxaCode 创建小程序码
func (w *implWxmp) CreateLimitedWxaCode(appId, appSecret, path string, width int, args ...Param) ([]byte, error) {
	accessToken, err := w.GetAccessToken(appId, appSecret)
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &typwx.LimitedWxaCodeParam{
		Path: path,
		CommonWxaCodeParam: &typwx.CommonWxaCodeParam{
			EnvVersion: "release",
			Width:      width,
			AutoColor:  true,
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

	// 如果不是图像数据，那就是json错误数据
	if !hdutils.IsImageData(resp.Body()) {
		return nil, errors.New(hdutils.BytesToString(resp.Body()))
	}

	return resp.Body(), nil
}

// CreateUnLimitedWxaCode 创建小程序码
func (w *implWxmp) CreateUnLimitedWxaCode(appId, appSecret, scene, page string, width int, args ...Param) ([]byte, error) {
	accessToken, err := w.GetAccessToken(appId, appSecret)
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &typwx.UnLimitedWxaCodeParam{
		Scene: scene,
		Page:  page,
		CommonWxaCodeParam: &typwx.CommonWxaCodeParam{
			EnvVersion: "release",
			Width:      width,
			AutoColor:  true,
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

	// 如果不是图像数据，那就是json错误数据
	if !hdutils.IsImageData(resp.Body()) {
		return nil, errors.New(hdutils.BytesToString(resp.Body()))
	}

	return resp.Body(), nil
}

func Trial(param *typwx.CommonWxaCodeParam) {
	param.EnvVersion = "trial"
}

func Develop(param *typwx.CommonWxaCodeParam) {
	param.EnvVersion = "develop"
}
