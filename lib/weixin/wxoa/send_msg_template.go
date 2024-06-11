package wxoa

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type SendMessageTemplateArgument struct {
	AppId        string
	AppSecret    string
	ToUserOpenId string
	TemplateId   string
	JumpUrl      string
	MiniProgram  *SendMessageTemplateMiniProgram
}

type SendMessageTemplateMiniProgram struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

type templateSendMessageImpl struct {
	httpClient *resty.Client
	arg        *SendMessageTemplateArgument
}

type sendMessageTemplate struct {
	ToUser      string                          `json:"touser"`
	TemplateId  string                          `json:"template_id"`
	Url         string                          `json:"JumpUrl"`
	MiniProgram *SendMessageTemplateMiniProgram `json:"miniprogram"`
	Data        any                             `json:"data"`
}

type sendMessageTemplateLine struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

const (
	urlSendTemplateMessage = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
	defaultColor           = "#173177"
)

func NewTemplateSendMessage(arg *SendMessageTemplateArgument) (SendMessager, error) {
	var httpClient = resty.New()
	httpClient.SetTimeout(3 * time.Second)

	return &templateSendMessageImpl{
		httpClient: httpClient,
		arg:        arg,
	}, nil
}

// Send 发送模板消息
func (m templateSendMessageImpl) Send(kvs ...string) error {
	accessToken, err := New(m.arg.AppId, m.arg.AppSecret).GetAccessToken()
	if err != nil {
		return err
	}

	realMsg, err := m.getTemplateMessage(kvs...)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlSendTemplateMessage, accessToken)
	resp, err := m.httpClient.SetHeader("Content-Type", "application/json; charset=UTF-8").R().SetBody(realMsg).Post(url)
	if err != nil {
		return errors.Wrapf(err, "send template message, JumpUrl: %s, content: %v", url, kvs)
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.Wrapf(err, "send template message, JumpUrl: %s, content: %v, statusCode: %d", url, kvs, resp.StatusCode())
	}

	return nil
}

func (m templateSendMessageImpl) getTemplateMessage(kvs ...string) (*sendMessageTemplate, error) {
	if len(kvs)%2 == 1 {
		return nil, errors.New("invalid key value content")
	}

	data := make(map[string]*sendMessageTemplateLine)
	for i := 0; i < len(kvs); i += 2 {
		data[kvs[i]] = &sendMessageTemplateLine{
			Value: kvs[i+1],
			Color: defaultColor,
		}
	}

	msg := &sendMessageTemplate{
		ToUser:      m.arg.ToUserOpenId,
		TemplateId:  m.arg.TemplateId,
		Url:         m.arg.JumpUrl,
		MiniProgram: m.arg.MiniProgram,
		Data:        data,
	}
	return msg, nil
}
