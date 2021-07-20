package alidts

import (
	"github.com/hamba/avro"
	"github.com/mitchellh/mapstructure"
)

type AliDts struct {
	schema avro.Schema
}

func New() (*AliDts, error) {
	s, err := avro.Parse(ALIYUN_DTS_SCHEMA)
	if err != nil {
		return nil, err
	}

	return &AliDts{
		schema: s,
	}, nil
}

// GetRecord 获取DTS的消息记录
func (ad *AliDts) GetRecord(data []byte) (*DtsRecord, error) {
	var v interface{}

	err := avro.Unmarshal(ad.schema, data, &v)
	if err != nil {
		return nil, err
	}

	var r DtsRecord
	err = mapstructure.Decode(v, &r)
	if err != nil {
		return nil, err
	}

	err = r.parse()
	if err != nil {
		return nil, err
	}

	return &r, nil
}
