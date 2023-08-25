package hdsdk

import (
	"bytes"
	"fmt"
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
	v                 *viper.Viper
	app               string
	env               string
	envPrefix         string          // 环境变量前缀
	rootParts         []string        // 根目录的部分，例如/setting/app, 则为[]string{"config", "app"}
	configType        string          // 配置内容类型，e,g: toml, json
	disableRemoteEnvs []string        // 禁用remote配置的环境列表
	fileOption        *OptionFile     // 文件配置选项
	remoteOptions     []*OptionRemote // 远程加载配置选项
	watchOption       *OptionWatch    // 检测配置变化选项
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
	defaultValue = struct {
		EnvPrefix         string
		RootParts         []string
		ConfigType        string
		FileOption        *OptionFile
		RemoteUrlKey      string
		RemoteProvider    string
		RemoteUrl         string
		WatchOption       *OptionWatch
		DisableRemoteEnvs []string
	}{
		EnvPrefix:  "HD",
		RootParts:  []string{"setting", "app"}, // 缺省的Root路径Parts, 这里定义为[]string，方便用path.Join或者filepath.Join
		ConfigType: "toml",                     // 缺省的配置文件类型
		FileOption: &OptionFile{
			dirs: make([]string, 0),
		},
		RemoteUrlKey:   "sdk.etcd.url",          // 默认etcd在文件中定义的key
		RemoteProvider: "etcd3",                 // 默认的remote provider
		RemoteUrl:      "http://127.0.0.1:2379", // 默认的remote url
		WatchOption: &OptionWatch{
			enabled:     true, // 默认是否检测远程配置变更
			effectDelay: 5,    // 配置变化生效时间为5秒
		},
		DisableRemoteEnvs: []string{"", "local"}, // 默认无环境或者local环境不需要加载remote配置
	}
)

const (
	// 最小化的配置,保证日志工作正常
	tplMinimalConfigContent = `
	[sdk]
	   [sdk.log]
	       level = "debug"
	       filename = "%s.log"
	       [sdk.log.rotate]
	           max_age = 168
	           rotation_time=24`
)

// NewConfig args[0]为env
func NewConfig(app, env string, options ...ConfigOption) *ViperConfig {
	c := &ViperConfig{
		v:                 viper.New(),
		app:               app,
		env:               env,
		envPrefix:         defaultValue.EnvPrefix,
		configType:        defaultValue.ConfigType,
		disableRemoteEnvs: defaultValue.DisableRemoteEnvs, // 禁用remote配置加载的环境列表
		fileOption:        defaultValue.FileOption,
		watchOption:       defaultValue.WatchOption,
	}

	for _, option := range options {
		option(c)
	}

	// 必须设置config的类型
	c.v.SetConfigType(c.configType)

	return c
}

func (c *ViperConfig) Load(configVar any) error {
	err := c.LoadLocal(nil)
	if err != nil {
		return err
	}

	// 只有当前环境不在disable列表时才需要加载remote配置
	if !utils.Contains(c.disableRemoteEnvs, c.env) {
		return c.LoadRemote(configVar)
	}

	return c.unmarshal(configVar)
}

// LoadLocal 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中）
// - env: 环境变量配置(高)
func (c *ViperConfig) LoadLocal(configVar any) error {
	// 如果环境变量为空，则加载最小基本配置
	if c.env == "" {
		c.loadMinimal()
		return nil
	}

	// 尝试从环境变量中获取配置信息
	c.loadFromEnv()

	// 尝试从配置文件中获取配置信息
	err := c.loadFromFile()
	if err != nil {
		return errors.Wrap(err, "load config from file")
	}

	return c.unmarshal(configVar)
}

func (c *ViperConfig) unmarshal(configVar any) error {
	// 尝试unmarshal所有配置数据
	if configVar != nil && len(c.v.AllKeys()) > 0 {
		err := c.v.Unmarshal(configVar)
		if err != nil {
			return errors.Wrap(err, "unmarshal config")
		}
	}
	return nil
}

func (c *ViperConfig) LoadRemote(configVar any) error {
	// 尝试从远程配置信息
	err := c.loadFromRemote()
	if err != nil {
		utils.LogError("load config from remote", "err", err)
	} else {
		// 如果加载remote成功，则尝试监控配置变化
		if c.watchOption.enabled {
			err = c.watchRemote(configVar)
			if err != nil {
				utils.LogError("watch remote config change", "err", err)
			}
		}
	}

	return c.unmarshal(configVar)
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

			err = c.unmarshal(configVar)
			if err != nil {
				utils.LogError("unable to unmarshal remote config", "err", err)
			}
		}
	}()
	return nil
}

func WithConfigFile(filepath string) ConfigOption {
	return func(c *ViperConfig) {
		c.fileOption.configFile = filepath
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

func WithDisableRemoteEnvs(args ...string) ConfigOption {
	return func(c *ViperConfig) {
		disableRemoteEnvs := make([]string, 0)
		c.disableRemoteEnvs = append(disableRemoteEnvs, args...)
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

func (c *ViperConfig) loadMinimal() {
	_ = c.ReadString(fmt.Sprintf(tplMinimalConfigContent, c.app))
}

func (c *ViperConfig) loadFromFile() error {
	// 如果指定了配置文件
	if c.fileOption.configFile != "" {
		c.v.SetConfigFile(c.fileOption.configFile)
	} else { //未指定在当前目录和父级目录找
		// 未指定搜索路径，使用缺省值setting/app/<app>
		if len(c.fileOption.dirs) == 0 {
			c.fileOption.dirs = append(c.fileOption.dirs,
				filepath.Join(filepath.Join(defaultValue.RootParts...), c.app),
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
// windows下测试: e,g: type test.txt | etcdctl.exe put /setting/app/hello/test
func (c *ViperConfig) loadFromRemote() error {
	if len(c.remoteOptions) == 0 {
		c.remoteOptions = append(c.remoteOptions, c.getDefaultRemoteOption())
	}

	for _, option := range c.remoteOptions {
		err := c.v.AddRemoteProvider(option.provider, option.url, option.path)
		if err != nil {
			return errors.Wrapf(err, "add remote provider, provider: %s, url: %s, path: %s", option.provider, option.url, option.path)
		}
	}

	err := c.v.ReadRemoteConfig()
	if err != nil {
		return errors.Wrapf(err, "read remote config")
	}

	for _, option := range c.remoteOptions {
		utils.LogDebug("load config from remote", "provider", option.provider, "url", option.url, "path", option.path)
	}
	return nil
}

func (c *ViperConfig) getDefaultRemoteOption() *OptionRemote {
	// 加载远程配置的时候优先从之前已经读取的配置，例如文件配置中取remoteUrl
	url := defaultValue.RemoteUrl
	if v := c.v.GetString(defaultValue.RemoteUrlKey); v != "" {
		url = v
	}

	return &OptionRemote{
		provider: defaultValue.RemoteProvider,
		url:      url,
		path:     path.Join("/", path.Join(defaultValue.RootParts...), c.app), // 具体app的具体环境的配置保存在该路径下： /setting/app/<app>

	}
}
