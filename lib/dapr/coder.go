package dapr

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/sqids/sqids-go"
	"io"
	"sync"
)

type EncoderDecoder interface {
	Encode(ids ...uint64) (string, error)
	DecodeUint64(code string) (uint64, error)
	DecodeUint64Slice(code string) ([]uint64, error)
}

type coderImpl struct {
	sqids  *sqids.Sqids
	secret []byte
}

const (
	saltLength = 4
	codeLength = 8
)

var (
	_onceCoder sync.Once
	_coder     *coderImpl
)

func Coder(secret []byte) EncoderDecoder {
	_onceCoder.Do(func() {
		s, _ := sqids.New(sqids.Options{
			MinLength: codeLength,
		})
		_coder = &coderImpl{
			sqids:  s,
			secret: secret,
		}
	})
	return _coder
}

func (impl coderImpl) Encode(ids ...uint64) (string, error) {
	s, err := impl.sqids.Encode(ids)
	if err != nil {
		return "", err
	}
	return impl.encodeWithSalt(s)
}

func (impl coderImpl) DecodeUint64(ciphertext string) (uint64, error) {
	plainText, err := impl.decodeWithSalt(ciphertext)
	if err != nil {
		return 0, err
	}

	uint64s := impl.sqids.Decode(plainText)
	if len(uint64s) <= 0 {
		return 0, errors.New("invalid code")
	}
	return uint64s[0], nil
}

func (impl coderImpl) DecodeUint64Slice(ciphertext string) ([]uint64, error) {
	plainText, err := impl.decodeWithSalt(ciphertext)
	if err != nil {
		return nil, err
	}

	return impl.sqids.Decode(plainText), nil
}

// 加密函数
func (impl coderImpl) encodeWithSalt(plaintext string) (string, error) {
	salt := make([]byte, saltLength) // 生成8字节的盐
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	// 将盐和明文拼接在一起
	data := append(impl.secret, []byte(plaintext)...)
	// 使用固定密钥进行XOR加密
	for i := range data {
		data[i] ^= impl.secret[i%len(impl.secret)]
	}
	// 返回Base64编码后的结果
	return base64.StdEncoding.EncodeToString(data), nil
}

// 解密函数
func (impl coderImpl) decodeWithSalt(ciphertextBase64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	// 使用固定密钥进行XOR解密
	for i := range data {
		data[i] ^= impl.secret[i%len(impl.secret)]
	}

	// 检查并移除盐
	if len(data) < saltLength {
		return "", errors.New("ciphertext too short")
	}
	plaintext := data[saltLength:]

	return string(plaintext), nil
}
