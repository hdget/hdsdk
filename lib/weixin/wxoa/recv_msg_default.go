package wxoa

import (
	"encoding/xml"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
)

type defaultRecvMessageImpl struct {
	*RecvMessage
}

func NewDefaultRecvMessage(data []byte) (RecvMessager, error) {
	var msg RecvMessage
	err := xml.Unmarshal(data, &msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal default message, data: %s", convert.BytesToString(data))
	}

	return &defaultRecvMessageImpl{RecvMessage: &msg}, nil
}

func (m *defaultRecvMessageImpl) Handle() ([]byte, error) {
	return m.ReplyText("")
}
