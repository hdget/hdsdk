package wx

import (
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/hdget/hdsdk/lib/wx/wxmp"
	"github.com/hdget/hdsdk/lib/wx/wxoa"
	"github.com/hdget/hdsdk/lib/wx/wxqr"
	"github.com/pkg/errors"
)

var (
	_wxmp = wxmp.New()
	_wxqr = wxqr.New()
	_wxoa = wxoa.New()
)

// WxmpAuth 小程序凭证校验, 返回openid, unionid
func WxmpAuth(appId, appSecret, code string) (string, string, error) {
	session, err := _wxmp.Auth(appId, appSecret, code)
	if err != nil {
		return "", "", errors.Wrap(err, "auth credential")
	}

	return session.OpenId, session.UnionId, nil
}

// WxmpDecryptUserInfo 微信小程序解密用户信息
func WxmpDecryptUserInfo(appId, encryptedData, iv string) (*typwx.WxmpUserInfo, error) {
	return _wxmp.DecryptUserInfo(appId, encryptedData, iv)
}

// WxmpDecryptMobileInfo 微信小程序解密手机信息
func WxmpDecryptMobileInfo(appId, encryptedData, iv string) (*typwx.WxmpMobileInfo, error) {
	return _wxmp.DecryptMobileInfo(appId, encryptedData, iv)
}

// WxmpGetUserPhoneNumber 新版获取用户手机号码
func WxmpGetUserPhoneNumber(appId, appSecret, code string) (*typwx.WxmpMobileInfo, error) {
	return _wxmp.GetUserPhoneNumber(appId, appSecret, code)
}

func WxmpCreateLimitedWxaCode(appId, appSecret, path string, width int, args ...wxmp.Param) ([]byte, error) {
	return _wxmp.CreateLimitedWxaCode(appId, appSecret, path, width, args...)
}

func WxmpCreateUnLimitedWxaCode(appId, appSecret, scene, page string, width int, args ...wxmp.Param) ([]byte, error) {
	return _wxmp.CreateUnLimitedWxaCode(appId, appSecret, scene, page, width, args...)
}

func WxmpGetAccessToken(appId, appSecret string) (string, error) {
	return _wxmp.GetAccessToken(appId, appSecret)
}

// WxqrGetWxId 微信二维码扫码登录
func WxqrGetWxId(appId, appSecret, code string) (string, string, error) {
	return _wxqr.GetWxId(appId, appSecret, code)
}

// WxoaGetSignature 微信公众号获取签名
func WxoaGetSignature(appId, appSecret, url string) (*typwx.WxoaSignature, error) {
	return _wxoa.GetSignature(appId, appSecret, url)
}
