package wxoa

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"time"
)

type RecvMessager interface {
	Handle() ([]byte, error)
	ReplyText(content string) ([]byte, error)
}

type MessageKind int

const (
	MessageKindUnknown MessageKind = iota
	MessageKindEvent
)

func GetMessageKind(data []byte) (MessageKind, error) {
	var msg RecvMessage
	err := xml.Unmarshal(data, &msg)
	if err != nil {
		return MessageKindUnknown, err
	}

	msgKind := MessageKindUnknown
	switch msg.MsgType {
	case "event":
		msgKind = MessageKindEvent
	}
	return msgKind, nil
}

func (m *RecvMessage) ReplyText(content string) ([]byte, error) {
	reply := SendMessageText{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		MsgId:        m.MsgId,
		Content:      content,
	}

	output, err := xml.MarshalIndent(reply, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal text message, reply: %v", reply)
	}

	return output, nil
}
