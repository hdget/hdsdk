package wxoa

import "encoding/xml"

// Signature signature接口返回结果
type Signature struct {
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"Url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// Ticket 类型
type Ticket struct {
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

type RecvMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	MsgId        int64
}

type RecvEventMessage struct {
	*RecvMessage
	Event    string
	EventKey string
	Ticket   string
}

type SendMessageText struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
}
