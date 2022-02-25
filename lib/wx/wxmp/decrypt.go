package wxmp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/pkg/errors"
)

var (
	ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

func (w *implWxmp) DecryptUserInfo(appId, encryptedData, iv string) (*typwx.WxmpUserInfo, error) {
	sessKey, err := _cache.GetSessKey(appId)
	if err != nil {
		return nil, errors.Wrap(err, "session key not found, you should invoke wx.login() firstly")
	}

	cipherText, err := decrypt(appId, sessKey, encryptedData, iv)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt encrypted data")
	}

	var userInfo typwx.WxmpUserInfo
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal to WxmpUserInfo")
	}

	if userInfo.Watermark.AppId != appId {
		return nil, ErrAppIDNotMatch
	}

	return &userInfo, nil
}

func (w *implWxmp) DecryptMobileInfo(appId, encryptedData, iv string) (*typwx.WxmpMobileInfo, error) {
	sessKey, err := _cache.GetSessKey(appId)
	if err != nil {
		return nil, errors.Wrap(err, "session key not found, you should invoke wx.login() firstly")
	}

	cipherText, err := decrypt(appId, sessKey, encryptedData, iv)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt encrypted data")
	}

	var mobileInfo typwx.WxmpMobileInfo
	err = json.Unmarshal(cipherText, &mobileInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal to WxmpUserInfo")
	}

	return &mobileInfo, nil
}

// 解密加密信息获取微信用户信息
func decrypt(appId, sessionKey, encryptedData, iv string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}
