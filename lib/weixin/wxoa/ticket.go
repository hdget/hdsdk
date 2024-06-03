package wxoa

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func (impl *wxoaImpl) getTicket() (string, error) {
	cachedTicket, err := impl.Cache.GetTicket()
	if err != nil {
		return "", errors.Wrap(err, "get wxoa cached ticket")
	}
	if cachedTicket != "" {
		return cachedTicket, nil
	}

	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return "", err
	}

	wxTicket, err := impl.requestTicket(accessToken)
	if err != nil {
		return "", err
	}

	// 忽略保存ticket错误
	err = impl.Cache.SetTicket(wxTicket.Value, wxTicket.ExpiresIn)
	if err != nil {
		return "", errors.Wrap(err, "set wxoa ticket to cache")
	}

	return wxTicket.Value, nil
}

// requestTicket jssdk获取凭证
func (impl *wxoaImpl) requestTicket(accessToken string) (*Ticket, error) {
	wxUserTicketTmpl := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	wxUserTicketURL := fmt.Sprintf(wxUserTicketTmpl, accessToken)
	client := resty.New()
	resp, err := client.R().Get(wxUserTicketURL)
	if err != nil {
		return nil, err
	}

	var ticket Ticket
	err = impl.ParseResult(resp.Body(), &ticket)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal weixin api response")
	}

	return &ticket, nil
}
