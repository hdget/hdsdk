package zerolog

import (
	"fmt"
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"os"
	"path/filepath"
	"runtime"
)

type zerologProviderConfig struct {
	Rotate   *rotateConfig `mapstructure:"rotate"`   // 日志文件截断的设置
	Dir      string        `mapstructure:"dir"`      // 日志目录
	Filename string        `mapstructure:"filename"` // 日志文件名
	Level    string        `mapstructure:"level"`    // 默认日志级别
}

type rotateConfig struct {
	MaxAge    int  `mapstructure:"max_age"`    // 多少天以后的日志删除
	MaxBackup int  `mapstructure:"max_backup"` // 保留多少个日志文件
	MaxSize   int  `mapstructure:"max_size"`   // 日志文件为多大开始rotate
	Compress  bool `mapstructure:"compress"`   // 是否压缩日志文件
}

const (
	configSection      = "sdk.log"
	linuxDefaultDir    = "/var/log"
	nonLinuxDefaultDir = "logs"
)

var (
	defaultConfig = &zerologProviderConfig{
		Rotate: &rotateConfig{
			MaxAge: 7, // 最大保存时间7天(单位hour)
		},
		Dir:      "logs",
		Filename: "app.log",
		Level:    "debug",
	}
)

// NewConfig 解析Config
func newConfig(configProvider intf.ConfigProvider) (*zerologProviderConfig, error) {
	if configProvider == nil {
		return getDefaultConfig(), nil
	}

	var c *zerologProviderConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return getDefaultConfig(), nil
	}

	if c == nil {
		return getDefaultConfig(), nil
	}

	// validate sdkConfig
	if c.Filename == "" || c.Rotate == nil {
		return nil, errdef.ErrInvalidConfig
	}

	return c, nil
}

func getDefaultConfig() *zerologProviderConfig {
	if dir, err := os.Getwd(); err == nil {
		guessAppName := filepath.Base(dir)
		defaultConfig.Filename = fmt.Sprintf("%s.log", guessAppName)
	}

	switch runtime.GOOS {
	case "linux":
		defaultConfig.Dir = linuxDefaultDir
	default:
		defaultConfig.Dir = nonLinuxDefaultDir
	}
	return defaultConfig
}
