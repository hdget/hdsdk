package wxoa

import (
	"crypto/sha1"
	"fmt"
	"github.com/hdget/hdsdk/v2/lib/weixin/base"
	"github.com/hdget/hdsdk/v2/lib/weixin/types"
	"github.com/hdget/hdutils/hash"
	"github.com/pkg/errors"
	"time"
)

type ApiWxoa interface {
	GetSignature(url string) (*Signature, error) // jsapi_ticket获取签名
	GetAccessToken() (string, error)
}

type wxoaImpl struct {
	*base.ApiWeixin
}

var (
	_ ApiWxoa = (*wxoaImpl)(nil)
)

func New(appId, appSecret string) ApiWxoa {
	return &wxoaImpl{
		ApiWeixin: base.New(types.WeixinAppWxoa, appId, appSecret),
	}
}

// nolint:recheck
func (impl *wxoaImpl) GetSignature(url string) (*Signature, error) {
	ticket, err := impl.getTicket()
	if err != nil {
		return nil, err
	}

	signature, err := impl.generateSignature(ticket, url)
	if err != nil {
		return nil, err
	}

	if signature == nil || signature.Signature == "" {
		return nil, errors.New("invalid signature")
	}

	return signature, nil
}

// 生成微信签名
func (impl *wxoaImpl) generateSignature(ticket, url string) (*Signature, error) {
	now := time.Now().Unix()
	nonceStr := hash.GenerateRandString(32)
	s := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket,
		nonceStr,
		now,
		url,
	)

	// 获取signature
	h := sha1.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return nil, err
	}
	hashValue := fmt.Sprintf("%x", h.Sum(nil))

	return &Signature{
		AppID:     impl.AppId,
		Ticket:    ticket,
		Noncestr:  nonceStr,
		Url:       url,
		Timestamp: now,
		Signature: hashValue,
	}, nil
}
