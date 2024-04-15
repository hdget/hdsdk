// Package config
// the default setting hierarchy looks like below:
//
//	...
//	setting/app/<app>/<app>.test.toml
//	setting/dapr/*
//	...
package viper

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/hdget/hdsdk/v2/types"
	"github.com/hdget/hdutils"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

// viperConfigLoader 命令行配置
type viperConfigLoader struct {
	local      *viper.Viper
	arg        *types.ConfigArgument
	envPrefix  string      // 环境变量前缀
	baseDir    string      // 配置文件所在的BaseDir
	configType string      // 配置内容类型，e,g: toml, json
	fileOption *fileOption // 文件配置选项
	content    []byte      // 如果用WithConfigContent指定了配置内容，则这里不为空
}

// 缺省配置选项
var (
	defaultValue = struct {
		envPrefix  string
		baseDir    string
		configType string
		fileOption *fileOption
	}{
		envPrefix:  "HD",
		baseDir:    filepath.Join("setting", "app"), // 其他环境的BaseDir
		configType: "toml",                          // 缺省的配置文件类型
		fileOption: &fileOption{
			dirs: make([]string, 0),
		},
	}
)

const (
	// 最小化的配置,保证日志工作正常
	tplMinimalConfigContent = `
[sdk]
  [sdk.logger]
	   level = "debug"
	   filename = "%s.log"
	   [sdk.logger.rotate]
		   max_age = 7`
)

// New 初始化config provider
func New(arg *types.ConfigArgument, options ...Option) (intf.ConfigProvider, error) {
	provider := &viperConfigLoader{
		local:      viper.New(),
		arg:        arg,
		envPrefix:  defaultValue.envPrefix,
		baseDir:    defaultValue.baseDir,
		configType: defaultValue.configType,
		fileOption: defaultValue.fileOption,
	}

	for _, option := range options {
		option(provider)
	}

	err := provider.loadLocal()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (p *viperConfigLoader) Unmarshal(configVar any, args ...string) error {
	if len(args) > 0 {
		return p.local.UnmarshalKey(args[0], configVar)
	}
	return p.local.Unmarshal(configVar)
}

// ///////////////////////////////////////////////////////////////
// private functions
// //////////////////////////////////////////////////////////////

// Load 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中）
// - env: 环境变量配置(高)
func (p *viperConfigLoader) loadLocal() error {
	// 必须设置config的类型
	p.local.SetConfigType(p.configType)

	// 如果指定了配置内容，则合并
	if p.content != nil {
		_ = p.local.MergeConfig(bytes.NewReader(p.content))
	}

	// 如果环境变量为空，则加载最小基本配置
	if p.arg.Env == "" {
		return p.loadMinimal()
	}

	// 尝试从环境变量中获取配置信息
	p.loadFromEnv()

	// 尝试从配置文件中获取配置信息
	return p.loadFromFile()
}

// loadFromEnv 从环境文件中读取配置信息
func (p *viperConfigLoader) loadFromEnv() {
	// 如果设置了环境变量前缀，则尝试自动获取环境变量中的配置
	if p.envPrefix != "" {
		p.local.SetEnvPrefix(p.envPrefix)
		p.local.AutomaticEnv()
	}
}

func (p *viperConfigLoader) loadMinimal() error {
	minimalConfig := fmt.Sprintf(tplMinimalConfigContent, p.arg.App)
	return p.local.MergeConfig(bytes.NewReader([]byte(minimalConfig)))
}

func (p *viperConfigLoader) loadFromFile() error {
	// 找配置文件
	p.setupConfigFile()

	// 读取配置文件
	err := p.local.ReadInConfig()
	if err != nil {
		return err
	}

	hdutils.LogDebug("load configer from file", "file", p.local.ConfigFileUsed())
	return nil
}

func (p *viperConfigLoader) setupConfigFile() {
	// 如果指定了配置文件
	if p.fileOption.configFile != "" {
		p.local.SetConfigFile(p.fileOption.configFile)
		return
	}

	// 未指定配置文件
	{
		// 获取config filename
		configFileName := p.fileOption.filename
		if configFileName == "" {
			configFileName = p.getDefaultConfigFilename()
		}

		// 获取config dirs
		configDirs := p.fileOption.dirs
		if len(configDirs) == 0 {
			foundDir := p.findConfigDir()
			if foundDir != "" {
				configDirs = append(configDirs, foundDir)
			} else {
				hdutils.LogFatal("no config dir found", "app", p.arg.App, "env", p.arg.Env)
			}
		}

		// 设置搜索选项
		for _, dir := range configDirs {
			p.local.AddConfigPath(dir) // 指定目录
		}
		p.local.SetConfigName(configFileName)
	}
}

// getDefaultConfigFilename 缺省的配置文件名: <app>.<env>
func (p *viperConfigLoader) getDefaultConfigFilename() string {
	return strings.Join([]string{p.arg.App, p.arg.Env}, ".")
}

// findConfigDirs 缺省的配置文件名: <app>.<env>
func (p *viperConfigLoader) findConfigDir() string {
	// parent dir name
	dirName := filepath.Join(p.baseDir, p.arg.App)

	// iter to root directory
	absStartPath, err := filepath.Abs(".")
	if err != nil {
		return ""
	}

	matchFile := fmt.Sprintf("%s.%s.%s", p.arg.App, p.arg.Env, p.configType)
	currPath := absStartPath
	for {
		s := filepath.Join(currPath, dirName, matchFile)
		matches, err := filepath.Glob(s)
		if err == nil && len(matches) > 0 {
			return filepath.Join(currPath, dirName)
		}

		// If we're already at the root, break the loop
		if currPath == "/" {
			break
		}

		// else, get parent dir
		currPath = filepath.Dir(currPath)
	}
	return ""
}
