package dapr

import (
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
)

func Reply(event *common.InvocationEvent, resp interface{}) (*common.Content, error) {
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &common.Content{
		ContentType: ContentTypeJson,
		Data:        data,
		DataTypeURL: event.DataTypeURL,
	}, nil
}
