package zerolog

import (
	"github.com/hdget/sdk/provider/log/typlog"
	"github.com/hdget/sdk/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type ConfigLog struct {
	Rotate   *RotateLogConf `mapstructure:"rotate"`   // 日志文件截断的设置
	Filename string         `mapstructure:"filename"` // 日志文件名
	Level    string         `mapstructure:"level"`    // 默认日志级别
}

type BaseLogProvider struct {
}

// @desc 解析Config
// @author Ryan Fan
// @param baseconf.Configer
// @return 第一个返回值为zerolog配置， 第二个返回值为error
func getLogConfig(configer types.Configer) (*ConfigLog, error) {
	v := configer.GetLogConfig()
	if v == nil {
		return nil, typlog.ErrEmptyLogConfig
	}

	values, ok := v.(map[string]interface{})
	if !ok {
		return nil, typlog.ErrInvalidLogConfig
	}

	var conf ConfigLog
	err := mapstructure.Decode(values, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode log config")
	}

	// validate config
	if conf.Filename == "" || conf.Rotate == nil {
		return nil, typlog.ErrInvalidLogConfig
	}

	return &conf, nil
}
