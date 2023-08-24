package hdsdk

import (
	"bytes"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ConfigOption func(config *ViperConfig)

// ViperConfig 命令行配置
type ViperConfig struct {
	v             *viper.Viper
	app           string
	env           string
	envPrefix     string          // 环境变量前缀
	rootParts     []string        // 根目录的部分，例如/setting/app, 则为[]string{"config", "app"}
	configType    string          // 配置内容类型，e,g: toml, json
	fileOption    *OptionFile     // 文件配置选项
	remoteOptions []*OptionRemote // 远程加载配置选项
	watchOption   *OptionWatch    // 检测配置变化选项
}

type OptionFile struct {
	configFile string   // 指定的配置文件
	dirs       []string // 如果未指定配置文件情况下，搜索的目录
	filename   string   // 如果未指定配置文件情况下，搜索的文件名，不需要文件后缀
}

type OptionRemote struct {
	provider string
	url      string
	path     string
}

type OptionWatch struct {
	enabled     bool // 是否开启监控
	effectDelay int  // 配置变化生效的间隔时间
}

// 缺省配置选项
var (
	defaultEnvPrefix  = "HD"
	defaultRootParts  = []string{"config", "app"} // 缺省的Root路径Parts, 这里定义为[]string，方便用path.Join或者filepath.Join
	defaultConfigType = "toml"                    // 缺省的配置文件类型
	defaultFileOption = &OptionFile{
		dirs: make([]string, 0),
	}
	defaultRemoteOption = &OptionRemote{
		provider: "etcd3",
		url:      "http://127.0.0.1:2379",
	}
	defaultWatchOption = &OptionWatch{
		enabled:     true,
		effectDelay: 5, // 配置变化生效时间为5秒
	}
)

const (
	// 最小化的配置,保证日志工作正常
	minimalConfigContent = `
	[sdk]
	   [sdk.log]
	       level = "debug"
	       filename = "app.log"
	       [sdk.log.rotate]
	           max_age = 168
	           rotation_time=24`
)

// NewConfig args[0]为env
func NewConfig(app, env string, options ...ConfigOption) *ViperConfig {
	c := &ViperConfig{
		v:           viper.New(),
		app:         app,
		env:         env,
		envPrefix:   defaultEnvPrefix,
		rootParts:   defaultRootParts,
		configType:  defaultConfigType,
		fileOption:  defaultFileOption,
		watchOption: defaultWatchOption,
	}

	for _, option := range options {
		option(c)
	}

	// 必须设置config的类型
	c.v.SetConfigType(c.configType)

	return c
}

// Load 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
// - default: 默认变量配置(最低)
// - remote: kvstore配置（低）
// - configFile: 文件配置(中）
// - env: 环境变量配置(高)
func (c *ViperConfig) Load(configVar any) error {
	// 尝试从环境变量中获取配置信息
	c.loadFromEnv()

	// 尝试从远程配置信息
	isRemoteOk := true
	err := c.loadFromRemote()
	if err != nil {
		isRemoteOk = false
		utils.LogError("load config from remote", "err", err)
	}

	// 尝试从配置文件中获取配置信息
	err = c.loadFromFile()
	if err != nil {
		return errors.Wrap(err, "load config from file")
	}

	// 尝试监控配置变化
	if c.watchOption.enabled && isRemoteOk {
		err = c.watchRemote(configVar)
		if err != nil {
			utils.LogError("watch remote config change", "err", err)
		}
	}

	// 尝试unmarshal所有配置数据
	if len(c.v.AllKeys()) > 0 {
		err = c.v.Unmarshal(configVar)
		if err != nil {
			return errors.Wrap(err, "unmarshal config")
		}
	}

	return nil
}

func (c *ViperConfig) Read(data []byte) *ViperConfig {
	_ = c.v.MergeConfig(bytes.NewReader(data))
	return c
}

func (c *ViperConfig) ReadString(s string) *ViperConfig {
	_ = c.Read(utils.StringToBytes(s))
	return c
}

func (c *ViperConfig) watchRemote(configVar any) error {
	// 如果无任何远程配置设置，忽略
	if len(c.remoteOptions) == 0 {
		return nil
	}

	// currently, only tested with etcd support
	err := c.v.WatchRemoteConfigOnChannel()
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(c.watchOption.effectDelay)) // delay after each request

			err = c.v.Unmarshal(configVar)
			if err != nil {
				utils.LogError("unable to unmarshal remote config", "err", err)
			}
		}
	}()
	return nil
}

func WithConfigFile(filename string) ConfigOption {
	return func(c *ViperConfig) {
		c.fileOption.configFile = filename
	}
}

func WithEnvPrefix(envPrefix string) ConfigOption {
	return func(c *ViperConfig) {
		c.envPrefix = envPrefix
	}
}

func WithConfigDir(args ...string) ConfigOption {
	return func(c *ViperConfig) {
		c.fileOption.dirs = append(c.fileOption.dirs, args...)
	}
}

func WithConfigFilename(filename string) ConfigOption {
	return func(c *ViperConfig) {
		if path.Ext(filename) != "" {
			utils.LogWarn("filename should not contains suffix", "filename", filename)
		}
		c.fileOption.filename = filename
	}
}

func WithConfigType(configType string) ConfigOption {
	return func(c *ViperConfig) {
		if !utils.Contains(viper.SupportedExts, configType) {
			utils.LogFatal("set config type", "supported", viper.SupportedExts, "err", viper.UnsupportedConfigError(configType))
		}
		c.configType = configType
	}
}

func WithRemote(provider, url, path string) ConfigOption {
	return func(c *ViperConfig) {
		if !utils.Contains(viper.SupportedRemoteProviders, provider) {
			utils.LogFatal("set remote config provider", "supported", viper.SupportedRemoteProviders, "err", viper.UnsupportedRemoteProviderError(provider))
		}

		if c.remoteOptions == nil {
			c.remoteOptions = make([]*OptionRemote, 0)
		}

		c.remoteOptions = append(c.remoteOptions, &OptionRemote{
			provider: provider,
			url:      url,
			path:     path,
		})
	}
}

func WithWatch(enabled bool, effectDelay int) ConfigOption {
	return func(c *ViperConfig) {
		c.watchOption.enabled = enabled
		c.watchOption.effectDelay = effectDelay
	}
}

func WithRoot(args ...string) ConfigOption {
	return func(c *ViperConfig) {
		rootParts := make([]string, 0)
		c.rootParts = append(rootParts, args...)
	}
}

// ///////////////////////////////////////////////////////////////
// private functions
// //////////////////////////////////////////////////////////////

// loadFromEnv 从环境文件中读取配置信息
func (c *ViperConfig) loadFromEnv() {
	// 如果设置了环境变量前缀，则尝试自动获取环境变量中的配置
	if c.envPrefix != "" {
		c.v.SetEnvPrefix(c.envPrefix)
		c.v.AutomaticEnv()
	}
}

func (c *ViperConfig) loadFromFile() error {
	// 如果未指定环境，则获取最小化的配置数据
	if c.env == "" {
		_ = c.ReadString(minimalConfigContent)
		return nil
	}

	// 如果指定了配置文件
	if c.fileOption.configFile != "" {
		c.v.SetConfigFile(c.fileOption.configFile)
	} else { //未指定在当前目录和父级目录找
		if len(c.fileOption.dirs) == 0 {
			c.fileOption.dirs = append(c.fileOption.dirs,
				filepath.Join(filepath.Join(defaultRootParts...), c.app), // config/app/<app>
			)
		}

		if c.fileOption.filename == "" {
			// 缺省的配置文件名: <app>.<env>
			c.fileOption.filename = strings.Join([]string{c.app, c.env}, ".")
		}

		// 先加入指定目录
		for _, dir := range c.fileOption.dirs {
			c.v.AddConfigPath(dir) // 指定目录
		}
		// 再默认添加指定目录的上级目录
		for _, dir := range c.fileOption.dirs {
			c.v.AddConfigPath(filepath.Join("..", dir))
		}
		c.v.SetConfigName(c.fileOption.filename)
	}

	err := c.v.ReadInConfig()
	if err != nil {
		// 如果配置文件找到，但读取时碰到其他问题需要中止
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		utils.LogDebug("load config from file", "file", c.v.ConfigFileUsed())
	}

	return nil
}

// loadFromRemote 尝试从远程kvstore中获取配置信息
func (c *ViperConfig) loadFromRemote() error {
	if len(c.remoteOptions) == 0 {
		// 缺省的配置路径: config/app/<app>/<env>
		defaultRemoteOption.path = path.Join("/", path.Join(defaultRootParts...), c.app, c.env)
		c.remoteOptions = []*OptionRemote{
			defaultRemoteOption,
		}
	}

	for _, option := range c.remoteOptions {
		err := c.v.AddRemoteProvider(option.provider, option.url, option.path)
		if err != nil {
			return errors.Wrapf(err, "add remote provider, provider: %s, url: %s, path: %s", option.provider, option.url, option.path)
		}
	}

	if len(c.remoteOptions) > 0 {
		err := c.v.ReadRemoteConfig()
		if err != nil {
			return errors.Wrapf(err, "read remote config")
		}
		utils.LogDebug("load config from remote")
	}

	return nil
}
