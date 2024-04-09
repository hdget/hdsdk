package viper

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdutils"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"path/filepath"
	"strings"
	"sync"
)

// viperConfigProvider 命令行配置
type viperConfigProvider struct {
	mu            sync.RWMutex
	local         *viper.Viper
	remote        *viper.Viper
	app           string
	env           string
	envPrefix     string          // 环境变量前缀
	rootParts     []string        // 根目录的部分，例如/setting/app, 则为[]string{"configer", "app"}
	configType    string          // 配置内容类型，e,g: toml, json
	fileOption    *fileOption     // 文件配置选项
	remoteOptions []*RemoteOption // 远程加载配置选项
	content       []byte          // 如果用WithConfigContent指定了配置内容，则这里不为空
}

const (
	sdkConfigSection = "sdk"
)

// 缺省配置选项
var (
	defaultValue = struct {
		EnvPrefix  string
		RootParts  []string
		ConfigType string
		FileOption *fileOption
	}{
		EnvPrefix:  "HD",
		RootParts:  []string{"setting", "app"}, // 缺省的Root路径Parts, 这里定义为[]string，方便用path.Join或者filepath.Join
		ConfigType: "toml",                     // 缺省的配置文件类型
		FileOption: &fileOption{
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

type Params struct {
	App     string   `name:"app"`
	Env     string   `name:"env"`
	Options []Option `name:"options"`
}

// New 初始化config provider
func New(app, env string, options ...Option) (intf.ConfigProvider, error) {
	provider := &viperConfigProvider{
		local:      viper.New(),
		remote:     viper.New(),
		app:        app,
		env:        env,
		envPrefix:  defaultValue.EnvPrefix,
		rootParts:  defaultValue.RootParts,
		configType: defaultValue.ConfigType,
		fileOption: defaultValue.FileOption,
	}

	for _, option := range options {
		option(provider)
	}

	err := provider.Init()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (vcLoader *viperConfigProvider) UnmarshalProviderConfig(rawVal any) error {
	return vcLoader.local.UnmarshalKey(sdkConfigSection, rawVal)
}

// Init 初始化configProvider, 加载本地配置到变量configVar
func (vcLoader *viperConfigProvider) Init(args ...any) error {
	return vcLoader.loadLocal()
}

// ///////////////////////////////////////////////////////////////
// private functions
// //////////////////////////////////////////////////////////////

// Load 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中）
// - env: 环境变量配置(高)
func (vcLoader *viperConfigProvider) loadLocal() error {
	// 必须设置config的类型
	vcLoader.local.SetConfigType(vcLoader.configType)

	// 如果指定了配置内容，则合并
	if vcLoader.content != nil {
		_ = vcLoader.local.MergeConfig(bytes.NewReader(vcLoader.content))
	}

	// 如果环境变量为空，则加载最小基本配置
	if vcLoader.env == "" {
		return vcLoader.loadMinimal()
	}

	// 尝试从环境变量中获取配置信息
	vcLoader.loadFromEnv()

	// 尝试从配置文件中获取配置信息
	return vcLoader.loadFromFile()
}

// loadFromEnv 从环境文件中读取配置信息
func (vcLoader *viperConfigProvider) loadFromEnv() {
	// 如果设置了环境变量前缀，则尝试自动获取环境变量中的配置
	if vcLoader.envPrefix != "" {
		vcLoader.local.SetEnvPrefix(vcLoader.envPrefix)
		vcLoader.local.AutomaticEnv()
	}
}

func (vcLoader *viperConfigProvider) loadMinimal() error {
	minimalConfig := fmt.Sprintf(tplMinimalConfigContent, vcLoader.app)
	return vcLoader.local.MergeConfig(bytes.NewReader(hdutils.StringToBytes(minimalConfig)))
}

func (vcLoader *viperConfigProvider) loadFromFile() error {
	// 如果指定了配置文件
	if vcLoader.fileOption.configFile != "" {
		vcLoader.local.SetConfigFile(vcLoader.fileOption.configFile)
	} else { //未指定在当前目录和父级目录找
		// 未指定搜索路径，使用缺省值setting/app/<app>
		if len(vcLoader.fileOption.dirs) == 0 {
			vcLoader.fileOption.dirs = append(vcLoader.fileOption.dirs,
				filepath.Join(filepath.Join(vcLoader.rootParts...), vcLoader.app),
			)
		}

		if vcLoader.fileOption.filename == "" {
			// 缺省的配置文件名: <app>.<env>
			vcLoader.fileOption.filename = strings.Join([]string{vcLoader.app, vcLoader.env}, ".")
		}

		// 先加入指定目录
		for _, dir := range vcLoader.fileOption.dirs {
			vcLoader.local.AddConfigPath(dir) // 指定目录
		}
		// 再默认添加指定目录的上级目录
		for _, dir := range vcLoader.fileOption.dirs {
			vcLoader.local.AddConfigPath(filepath.Join("..", dir))
		}
		vcLoader.local.SetConfigName(vcLoader.fileOption.filename)
	}

	err := vcLoader.local.ReadInConfig()
	if err != nil {
		// 如果配置文件找到，但读取时碰到其他问题需要中止
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		hdutils.LogDebug("load configer from file", "file", vcLoader.local.ConfigFileUsed())
	}

	return nil
}
