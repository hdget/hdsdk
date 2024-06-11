package wxoa

import (
	"encoding/xml"
	"github.com/hdget/hdutils/convert"
	"github.com/pkg/errors"
)

type eventRecvMessageImpl struct {
	*RecvEventMessage
	eventHandlers map[string]EventMessageHandler
}

type EventMessageHandler func(message *RecvEventMessage) ([]byte, error)

func NewEventRecvMessage(data []byte, eventHandlers map[string]EventMessageHandler) (RecvMessager, error) {
	var msg RecvEventMessage
	err := xml.Unmarshal(data, &msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal event message, data: %s", convert.BytesToString(data))
	}

	return &eventRecvMessageImpl{RecvEventMessage: &msg, eventHandlers: eventHandlers}, nil
}

func (m *eventRecvMessageImpl) Handle() ([]byte, error) {
	if handler, exists := m.eventHandlers[m.RecvEventMessage.Event]; exists {
		return handler(m.RecvEventMessage)
	}
	return nil, nil
}
