package wxoa

import (
	"encoding/json"
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
	Url          string // 公众号的跳转页面
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
	Url         string                          `json:"url"`
	MiniProgram *SendMessageTemplateMiniProgram `json:"miniprogram"`
	Data        any                             `json:"data"`
}

type sendMessageTemplateLine struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type sendResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Msgid   int    `json:"msgid"`
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
func (m templateSendMessageImpl) Send(contents map[string]string) error {
	accessToken, err := New(m.arg.AppId, m.arg.AppSecret).GetAccessToken()
	if err != nil {
		return err
	}

	msg, err := m.getTemplateMessage(contents)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(urlSendTemplateMessage, accessToken)
	resp, err := m.httpClient.SetHeader("Content-Type", "application/json; charset=UTF-8").R().SetBody(msg).Post(url)
	if err != nil {
		return errors.Wrapf(err, "send template message, url: %s, message: %v", url, msg)
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.Wrapf(err, "send template message, url: %s, message: %v, statusCode: %d", url, msg, resp.StatusCode())
	}

	var result sendResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return errors.Wrapf(err, "parse template message send result, url: %s, message: %v, statusCode: %d", url, msg, resp.StatusCode())
	}

	if result.Errcode != 0 {
		return errors.Wrapf(errors.New(result.Errmsg), "send template message, url: %s, message: %v, statusCode: %d", url, msg, resp.StatusCode())
	}

	return nil
}

func (m templateSendMessageImpl) getTemplateMessage(contents map[string]string) (*sendMessageTemplate, error) {
	data := make(map[string]*sendMessageTemplateLine)
	for k, v := range contents {
		data[k] = &sendMessageTemplateLine{
			Value: v,
			Color: defaultColor,
		}
	}

	msg := &sendMessageTemplate{
		ToUser:      m.arg.ToUserOpenId,
		TemplateId:  m.arg.TemplateId,
		Url:         m.arg.Url,
		MiniProgram: m.arg.MiniProgram,
		Data:        data,
	}
	return msg, nil
}
