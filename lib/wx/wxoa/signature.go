package wxoa

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/hdget/hdsdk/hdutils"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"time"
)

// nolint:recheck
func (w *implWxoa) GetSignature(appId, appSecret, url string) (*typwx.WxoaSignature, error) {
	ticket, err := w.getTicket(appId, appSecret)
	if err != nil {
		return nil, err
	}

	signature, err := genSignature(appId, ticket, url)
	if err != nil {
		return nil, err
	}

	if signature == nil || signature.Signature == "" {
		return nil, errors.New("invalid signature")
	}

	return signature, nil
}

// 生成微信签名
func genSignature(appId, ticket, url string) (*typwx.WxoaSignature, error) {
	now := time.Now().Unix()
	noncestr := hdutils.GenerateRandString(32)
	longstr := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket,
		noncestr,
		now,
		url,
	)

	// 获取signature
	h := sha1.New()
	_, err := h.Write([]byte(longstr))
	if err != nil {
		return nil, err
	}
	signature := fmt.Sprintf("%x", h.Sum(nil))

	return &typwx.WxoaSignature{
		AppID:     appId,
		Ticket:    ticket,
		Noncestr:  noncestr,
		Url:       url,
		Timestamp: now,
		Signature: signature,
	}, nil
}
