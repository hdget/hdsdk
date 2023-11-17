package zerolog

import (
	"fmt"
	"github.com/hdget/hdsdk/provider/log/typlog"
	"github.com/hdget/hdsdk/types"
	hdutils "github.com/hdget/hdutils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
)

type ConfigLog struct {
	Rotate   *RotateLogConf `mapstructure:"rotate"`   // 日志文件截断的设置
	Filename string         `mapstructure:"filename"` // 日志文件名
	Level    string         `mapstructure:"level"`    // 默认日志级别
}

type RotateLogConf struct {
	MaxAge       int    `mapstructure:"max_age"`
	RotationTime int    `mapstructure:"rotation_time"`
	Dirname      string `mapstructure:"dirname"` // 日志文件的保存目录名
	BaseDir      string `mapstructure:"basedir"` // 在linux环境下日志实际保存在<basedir>/<app>/<dirname>中,然后以link的方式创建dirname
}

const (
	defaultLinuxBaseDir = "/var/log"
)

var (
	defaultLogConfig = &ConfigLog{
		Rotate: &RotateLogConf{
			MaxAge:       168, // 最大保存时间7天(单位hour)
			RotationTime: 24,  // 日志切割时间间隔24小时（单位hour)
			BaseDir:      ".",
			Dirname:      "logs", // log directory name
		},
		Filename: "app.log",
		Level:    "debug",
	}
)

// getLogConfig 解析Config
func getLogConfig(configer types.Configer) (*ConfigLog, error) {
	if configer == nil {
		hdutils.LogWarn("empty configer")
		return getDefaultLogConfig(), nil
	}

	// if log config not found, use default one
	v := configer.GetLogConfig()
	if v == nil {
		hdutils.LogWarn("log config not found")
		return getDefaultLogConfig(), nil
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

func getDefaultLogConfig() *ConfigLog {
	// need to handle base dir and filename
	if dir, err := os.Getwd(); err == nil {
		guessAppName := filepath.Base(dir)
		defaultLogConfig.Filename = fmt.Sprintf("%s.log", guessAppName)
	}

	switch runtime.GOOS {
	case "linux":
		defaultLogConfig.Rotate.BaseDir = defaultLinuxBaseDir
	}
	return defaultLogConfig
}
