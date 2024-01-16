package zerolog

import (
	"fmt"
	"github.com/hdget/hdsdk/core/logger/errdef"
	"github.com/hdget/hdsdk/intf"
	"github.com/hdget/hdutils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
)

type zerologConfig struct {
	Rotate   *rotateConfig `mapstructure:"rotate"`   // 日志文件截断的设置
	Dir      string        `mapstructure:"dir"`      // 日志目录
	Filename string        `mapstructure:"filename"` // 日志文件名
	Level    string        `mapstructure:"level"`    // 默认日志级别
}

type rotateConfig struct {
	MaxAge    int  `mapstructure:"max_age"`  // 单位为天
	MaxBackup int  `mapstructure:"max_age"`  // 保留多少个日志文件
	MaxSize   int  `mapstructure:"max_size"` // 日志文件为多大开始rotate
	Compress  bool `mapstructure:"compress"` // 是否压缩日志文件
}

const (
	linuxDefaultDir    = "/var/log"
	nonLinuxDefaultDir = "logs"
)

var (
	defaultConfig = &zerologConfig{
		Rotate: &rotateConfig{
			MaxAge: 7, // 最大保存时间7天(单位hour)
		},
		Dir:      "logs",
		Filename: "app.logger",
		Level:    "debug",
	}
)

// NewConfig 解析Config
func NewConfig(configer intf.Configer) (*zerologConfig, error) {
	if configer == nil {
		hdutils.LogWarn("empty sdkConfig")
		return getDefaultConfig(), nil
	}

	// if logger sdkConfig not found, use default one
	v := configer.GetLogConfig()
	if v == nil {
		hdutils.LogWarn("logger sdkConfig not found")
		return getDefaultConfig(), nil
	}

	var conf zerologConfig
	err := mapstructure.Decode(v, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "decode logger configloader")
	}

	// validate sdkConfig
	if conf.Filename == "" || conf.Rotate == nil {
		return nil, errdef.ErrInvalidLogConfig
	}

	return &conf, nil
}

func getDefaultConfig() *zerologConfig {
	if dir, err := os.Getwd(); err == nil {
		guessAppName := filepath.Base(dir)
		defaultConfig.Filename = fmt.Sprintf("%s.logger", guessAppName)
	}

	switch runtime.GOOS {
	case "linux":
		defaultConfig.Dir = linuxDefaultDir
	default:
		defaultConfig.Dir = nonLinuxDefaultDir
	}
	return defaultConfig
}
