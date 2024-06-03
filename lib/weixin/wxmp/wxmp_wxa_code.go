package wxmp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdutils/cmp"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
)

type WxaCodeOption func(*WxaCodeConfig)

const (
	urlGetLimitedWxaCode   = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	urlGetUnlimitedWxaCode = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

// CreateLimitedWxaCode 创建小程序码
func (impl *wxmpImpl) CreateLimitedWxaCode(path string, width int, options ...WxaCodeOption) ([]byte, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &LimitedWxaCode{
		Path: path,
		WxaCodeConfig: &WxaCodeConfig{
			EnvVersion: "release",
			Width:      width,
			AutoColor:  true,
		},
	}
	for _, opt := range options {
		opt(body.WxaCodeConfig)
	}

	resp, err := resty.New().R().SetBody(body).Post(fmt.Sprintf(urlGetLimitedWxaCode, accessToken))
	if err != nil {
		return nil, err
	}

	// 如果不是图像数据，那就是json错误数据
	if !cmp.IsImageData(resp.Body()) {
		return nil, errors.New(convert.BytesToString(resp.Body()))
	}

	return resp.Body(), nil
}

// CreateUnLimitedWxaCode 创建小程序码
func (impl *wxmpImpl) CreateUnLimitedWxaCode(scene, page string, width int, options ...WxaCodeOption) ([]byte, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &UnlimitedWxaCode{
		Scene: scene,
		Page:  page,
		WxaCodeConfig: &WxaCodeConfig{
			EnvVersion: "release",
			Width:      width,
			AutoColor:  true,
		},
	}
	for _, opt := range options {
		opt(body.WxaCodeConfig)
	}

	resp, err := resty.New().R().SetBody(body).Post(fmt.Sprintf(urlGetUnlimitedWxaCode, accessToken))
	if err != nil {
		return nil, err
	}

	// 如果不是图像数据，那就是json错误数据
	if !cmp.IsImageData(resp.Body()) {
		return nil, errors.New(convert.BytesToString(resp.Body()))
	}

	return resp.Body(), nil
}

func Trial() WxaCodeOption {
	return func(c *WxaCodeConfig) {
		c.EnvVersion = "trial"
	}
}

func Develop() WxaCodeOption {
	return func(c *WxaCodeConfig) {
		c.EnvVersion = "develop"
	}
}
