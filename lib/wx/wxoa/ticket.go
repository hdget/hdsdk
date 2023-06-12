package wxoa

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/pkg/errors"
)

// WxoaTicket 类型
type WxoaTicket struct {
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

func (w *implWxoa) getTicket(appId, appSecret string) (string, error) {
	cachedTicket, err := _cache.GetTicket()
	if err != nil {
		return "", errors.Wrap(err, "get wxoa cached ticket")
	}
	if cachedTicket != "" {
		return cachedTicket, nil
	}

	accessToken, err := w.GetAccessToken(appId, appSecret)
	if err != nil {
		return "", err
	}

	wxTicket, err := w.requestTicket(accessToken)
	if err != nil {
		return "", err
	}

	// 忽略保存ticket错误
	err = _cache.SetTicket(wxTicket.Value, wxTicket.ExpiresIn)
	if err != nil {
		return "", errors.Wrap(err, "set wxoa ticket to cache")
	}

	return wxTicket.Value, nil
}

// requestTicket 二维码ticket
func (w *implWxoa) requestTicket(accessToken string) (*WxoaTicket, error) {
	wxUserTicketTmpl := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	wxUserTicketURL := fmt.Sprintf(wxUserTicketTmpl, accessToken)
	client := resty.New()
	resp, err := client.R().Get(wxUserTicketURL)
	if err != nil {
		return nil, err
	}

	// var ticket WxoaTicket
	ticket := &WxoaTicket{}
	err = json.Unmarshal(resp.Body(), ticket)
	if ticket.Value == "" {
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal to WxoaTicket")
		}

		// 如果unmarshal请求消息错误,尝试获取错误信息
		var errResp typwx.WxErrResponse
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			hdsdk.Logger.Error("unmarshal wx err response", "err", err)
		}
		return nil, errors.New(errResp.ErrMsg)
	}

	return ticket, nil
}
