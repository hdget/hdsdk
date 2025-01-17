package code

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
	"github.com/sqids/sqids-go"
	"io"
	"sync"
)

type Coder interface {
	Encrypt(secret []byte, ids ...int64) string
	DecryptInt64(secret []byte, ciphertext string) int64
	DecryptInt64Slice(secret []byte, ciphertext string) []int64
	Encode(ids ...int64) string
	DecodeInt64(code string) int64
	DecodeInt64Slice(code string) []int64
}

type coderImpl struct {
	sqids *sqids.Sqids
}

const (
	saltLength = 4
	codeLength = 8
)

var (
	_onceCoder sync.Once
	_coder     *coderImpl
)

func New() Coder {
	_onceCoder.Do(func() {
		s, _ := sqids.New(sqids.Options{
			MinLength: codeLength,
		})
		_coder = &coderImpl{
			sqids: s,
		}
	})
	return _coder
}

func (impl coderImpl) Encrypt(secret []byte, ids ...int64) string {
	uint64s := pie.Map(ids, func(v int64) uint64 { return uint64(v) })
	s, err := impl.sqids.Encode(uint64s)
	if err != nil {
		return ""
	}

	encoded, err := impl.encodeWithSalt(s, secret)
	if err != nil {
		return ""
	}
	return encoded
}

func (impl coderImpl) DecryptInt64(secret []byte, ciphertext string) int64 {
	plainText, err := impl.decodeWithSalt(ciphertext, secret)
	if err != nil {
		return 0
	}

	uint64s := impl.sqids.Decode(plainText)
	if len(uint64s) <= 0 {
		return 0
	}

	return int64(uint64s[0])
}

func (impl coderImpl) DecryptInt64Slice(secret []byte, ciphertext string) []int64 {
	plainText, err := impl.decodeWithSalt(ciphertext, secret)
	if err != nil {
		return nil
	}

	uint64s := impl.sqids.Decode(plainText)
	return pie.Map(uint64s, func(v uint64) int64 { return int64(v) })
}

func (impl coderImpl) Encode(ids ...int64) string {
	uint64s := pie.Map(ids, func(v int64) uint64 { return uint64(v) })
	s, err := impl.sqids.Encode(uint64s)
	if err != nil {
		return ""
	}
	return s
}

func (impl coderImpl) DecodeInt64(s string) int64 {
	uint64s := impl.sqids.Decode(s)
	if len(uint64s) <= 0 {
		return 0
	}

	return int64(uint64s[0])
}

func (impl coderImpl) DecodeInt64Slice(code string) []int64 {
	uint64s := impl.sqids.Decode(code)
	return pie.Map(uint64s, func(v uint64) int64 { return int64(v) })
}

// 加密函数
func (impl coderImpl) encodeWithSalt(plaintext string, secret []byte) (string, error) {
	salt := make([]byte, saltLength) // 生成8字节的盐
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	// 将盐和明文拼接在一起
	data := append(secret, []byte(plaintext)...)
	// 使用固定密钥进行XOR加密
	for i := range data {
		data[i] ^= secret[i%len(secret)]
	}
	// 返回Base64编码后的结果
	return base64.URLEncoding.EncodeToString(data), nil
}

// 解密函数
func (impl coderImpl) decodeWithSalt(ciphertextBase64 string, secret []byte) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	// 使用固定密钥进行XOR解密
	for i := range data {
		data[i] ^= secret[i%len(secret)]
	}

	// 检查并移除盐
	if len(data) < saltLength {
		return "", errors.New("ciphertext too short")
	}
	plaintext := data[saltLength:]

	return string(plaintext), nil
}
