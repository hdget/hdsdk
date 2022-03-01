package aliyun

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

// DtsRecord 原始的记录
type DtsRecord struct {
	Version         int                    `mapstructure:"version"`
	Id              int64                  `mapstructure:"id"`
	SourceTimeStamp int64                  `mapstructure:"sourceTimestamp"`
	SourceTxId      string                 `mapstructure:"sourceTxid"`
	ObjectName      map[string]string      `mapstructure:"objectName"` // 数据库名.表名
	Operation       string                 `mapstructure:"operation"`
	Fields          map[string]interface{} `mapstructure:"fields"`       // 字段slice
	BeforeImages    map[string]interface{} `mapstructure:"beforeImages"` // 改变前
	AfterImages     map[string]interface{} `mapstructure:"afterImages"`  // 改变后

	// 额外的字段
	Database    string
	Table       string
	TableFields []*DtsField
}

type DtsField struct {
	Name     string `mapstructure:"name"`
	DataType int    `mapstructure:"dataTypeNumber"`
}

type DtsFields struct {
	Items []*DtsField `mapstructure:"array"`
}

type DtsTypeTimestamp struct {
	Timestamp int64 `mapstructure:"timestamp"`
}

type DtsTypeDateTime struct {
	Year   map[string]interface{} `mapstructure:"year"`
	Month  map[string]interface{} `mapstructure:"month"`
	Day    map[string]interface{} `mapstructure:"day"`
	Hour   map[string]interface{} `mapstructure:"hour"`
	Minute map[string]interface{} `mapstructure:"minute"`
	Second map[string]interface{} `mapstructure:"second"`
}

type DtsTypeTimestampWithTimeZone struct {
	Value DtsTypeDateTime `mapstructure:"value"`
}

type DtsTypeValue struct {
	Value string `mapstructure:"value"`
}

type DtsTypeBytes struct {
	Value []byte `mapstructure:"value"`
}

func (r *DtsRecord) GetAfterColumns() map[string]string {
	return r.getColumns(r.AfterImages)
}

func (r *DtsRecord) GetBeforeColumns() map[string]string {
	return r.getColumns(r.BeforeImages)
}

// 解析一些东西
func (r *DtsRecord) parse() error {
	// 解析数据库名和表名
	if r.ObjectName != nil {
		name := r.ObjectName["string"]
		tokens := strings.Split(name, ".")
		switch len(tokens) {
		case 1:
			r.Database = tokens[0]
		case 2:
			r.Database = tokens[0]
			r.Table = tokens[1]
		}
	}

	return nil
}

// 获取ColValue
func (r *DtsRecord) getColValue(kv interface{}) string {
	if kv == nil {
		return ""
	}

	mapValues, ok := kv.(map[string]interface{})
	if !ok {
		return ""
	}

	ret := ""
	for k, v := range mapValues {
		switch k {
		case "com.alibaba.alidts.formats.avro.Character",
			"com.alibaba.alidts.formats.avro.BinaryGeometry":
			var vv DtsTypeBytes
			err := mapstructure.Decode(v, &vv)
			if err != nil {
				return ""
			}
			ret = string(vv.Value)
		case "com.alibaba.alidts.formats.avro.Integer",
			"com.alibaba.alidts.formats.avro.Decimal",
			"com.alibaba.alidts.formats.avro.Float",
			"com.alibaba.alidts.formats.avro.TextGeometry",
			"com.alibaba.alidts.formats.avro.TextObject":
			var vv DtsTypeValue
			err := mapstructure.Decode(v, &vv)
			if err != nil {
				return ""
			}
			ret = vv.Value
		case "com.alibaba.alidts.formats.avro.Timestamp":
			var vv DtsTypeTimestamp
			err := mapstructure.Decode(v, &vv)
			if err != nil {
				return ""
			}
			ret = strconv.FormatInt(vv.Timestamp, 10)
		case "com.alibaba.alidts.formats.avro.DateTime":
			var vv DtsTypeDateTime
			err := mapstructure.Decode(v, &vv)
			if err != nil {
				return ""
			}
			ret = fmt.Sprintf("%v-%v-%v %v:%v:%v",
				vv.Year["int"], vv.Month["int"], vv.Day["int"], vv.Hour["int"], vv.Minute["int"], vv.Second["int"])
		case "com.alibaba.alidts.formats.avro.TimestampWithTimeZone":
			var vv DtsTypeTimestampWithTimeZone
			err := mapstructure.Decode(v, &vv)
			if err != nil {
				return ""
			}
			ret = fmt.Sprintf("%v-%v-%v %v:%v:%v",
				vv.Value.Year["int"], vv.Value.Month["int"], vv.Value.Day["int"], vv.Value.Hour["int"], vv.Value.Minute["int"], vv.Value.Second["int"])
		}
	}

	return ret
}

func (r *DtsRecord) getColumns(images map[string]interface{}) map[string]string {
	imageArray := images["array"]
	if imageArray == nil {
		return nil
	}

	if len(r.TableFields) == 0 {
		fields, err := r.getFields()
		if err != nil {
			return nil
		}

		if len(fields.Items) == 0 {
			return nil
		}

		r.TableFields = fields.Items
	}

	array := imageArray.([]interface{})
	if len(r.TableFields) != len(array) {
		return nil
	}

	cols := make(map[string]string)
	for index, v := range array {
		field := r.TableFields[index]
		if field == nil {
			continue
		}

		cols[field.Name] = r.getColValue(v)
	}

	return cols
}

func (r *DtsRecord) getFields() (*DtsFields, error) {
	var fields DtsFields
	err := mapstructure.Decode(r.Fields, &fields)
	if err != nil {
		return nil, nil
	}

	return &fields, nil
}
