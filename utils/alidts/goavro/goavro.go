package goavro

import (
	"github.com/linkedin/goavro/v2"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type AliDts struct {
	codec *goavro.Codec
}

func New() (*AliDts, error) {
	c, err := goavro.NewCodec(ALIYUN_DTS_SCHEMA)
	if err != nil {
		return nil, err
	}

	return &AliDts{
		codec: c,
	}, nil
}

// GetRecord 获取DTS的消息记录
func (ad *AliDts) GetRecord(data []byte) (*DtsRecord, error) {
	// 将消息值转换成map[string]interface{} 嵌套的golang对象
	nativeObj, _, err := ad.codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	var r DtsRecord
	err = mapstructure.Decode(nativeObj, &r)
	if err != nil {
		return nil, err
	}

	// 解析数据库名和表名
	objectName := r.ObjectName["string"]
	if objectName != "" {
		tokens := strings.Split(objectName, ".")
		if len(tokens) == 2 {
			r.Database = tokens[0]
			r.Table = tokens[1]
		}
	}

	return &r, nil
}
