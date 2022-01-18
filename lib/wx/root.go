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

// 小程序凭证校验, 返回openid, unionid
func WxmpAuth(appId, appSecret, code string) (string, string, error) {
	session, err := _wxmp.Auth(appId, appSecret, code)
	if err != nil {
		return "", "", errors.Wrap(err, "auth credential")
	}

	return session.OpenId, session.UnionId, nil
}

// WxmpDecrypt 微信小程序解密
func WxmpDecrypt(appId, wxId, encryptedData, iv string) (*typwx.WxmpUserInfo, error) {
	return _wxmp.Decrypt(appId, wxId, encryptedData, iv)
}

// WxqrGetWxId 微信二维码扫码登录
func WxqrGetWxId(appId, appSecret, code string) (string, string, error) {
	return _wxqr.GetWxId(appId, appSecret, code)
}

// WxoaGetSignature 微信公众号获取签名
func WxoaGetSignature(appId, appSecret, url string) (*typwx.WxoaSignature, error) {
	return _wxoa.GetSignature(appId, appSecret, url)
}
