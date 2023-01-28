package hdsdk

import (
	"fmt"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"path"
	"strings"
)

// ViperConfig 命令行配置
type ViperConfig struct {
	v      *viper.Viper
	app    string
	env    string
	option *ConfigOption
}

type ConfigOption struct {
	Env  EnvOption  // 环境变量选项
	File FileOption // 配置文件选项
	Etcd EtcdOption // etcd选项
}

type EnvOption struct {
	Prefix string // 环境变量前缀
}

type FileOption struct {
	RootDir string // 配置文件所在的根目录
	BaseDir string // 配置文件的上级目录
	Suffix  string // .toml
}

type EtcdOption struct {
	Root string
	Url  string
}

// 缺省配置选项
var (
	defaultEnvOption = EnvOption{
		Prefix: "HDGET",
	}

	defaultFileOption = FileOption{
		RootDir: "setting",
		BaseDir: "app",
		Suffix:  "toml",
	}

	defaultEtcdOption = EtcdOption{
		Root: "setting",
		Url:  "http://127.0.0.1:2379",
	}
)

// NewConfig args[0]为env
func NewConfig(app string, args ...string) *ViperConfig {
	// 默认为dev
	var env string
	// 如果cli指定了env，则优先使用cli指定的env
	if len(args) > 0 {
		env = args[0]
	} else { // 尝试使用环境变量中使用的env
		env = os.Getenv(fmt.Sprintf("%s_ENV", defaultEnvOption.Prefix))
	}

	return &ViperConfig{
		v:   viper.New(),
		app: app,
		env: env,
		option: &ConfigOption{
			Env:  defaultEnvOption,
			File: defaultFileOption,
			Etcd: defaultEtcdOption,
		},
	}
}

// Load will load config from difference sources,
// the config will read with following precedence:
// - flag(命令行的设置是优先级最高的)
// - env(其次为环境变量)
// - config(再其次为配置文件)
// - key/value store(接着是KVStore(etcd)里面的设置)
// - default(最低的是默认配置)
func (c *ViperConfig) Load(args ...string) *viper.Viper {
	// 尝试从环境中读取配置
	c.loadFromEnv()

	// 尝试从远程KV store，例如etcd加载配置信息
	err := c.loadFromRemote()
	if err != nil {
		utils.LogWarn("load config from etcd", "err", err)
	}

	// 尝试从配置中读取配置信息
	err = c.loadFromFile(args...)
	if err != nil {
		utils.LogWarn("load config from file", "err", err)
	}

	return c.v
}

// SetFileRoot 设置配置文件的rootDir
func (c *ViperConfig) SetFileRoot(rootDir string) *ViperConfig {
	c.option.File.RootDir = rootDir
	return c
}

// SetFileBase 设置配置文件的baseDir
func (c *ViperConfig) SetFileBase(baseDir string) *ViperConfig {
	c.option.File.BaseDir = baseDir
	return c
}

// SetRemoteRoot 设置etcd的root
func (c *ViperConfig) SetRemoteRoot(rootDir string) *ViperConfig {
	c.option.Etcd.Root = rootDir
	return c
}

// SetRemoteUrl 设置远程provider的url
func (c *ViperConfig) SetRemoteUrl(url string) *ViperConfig {
	c.option.Etcd.Url = url
	return c
}

// ///////////////////////////////////////////////////////////////
// private functions
// //////////////////////////////////////////////////////////////

// loadFromEnv 从环境文件中读取配置信息
func (c *ViperConfig) loadFromEnv(args ...string) {
	// 获取有效的环境变量prefix, 默认为:
	envPrefix := defaultEnvOption.Prefix
	if len(args) > 0 && args[0] != "" {
		envPrefix = args[0]
	}
	c.v.SetEnvPrefix(envPrefix)
	c.v.AutomaticEnv()
}

func (c *ViperConfig) loadFromFile(args ...string) error {
	configFile := c.getConfigFile()
	if len(args) > 0 && strings.Trim(args[0], " ") != "" {
		configFile = args[0]
	}

	// optionally look for config in the working directory
	c.v.AddConfigPath(".")
	c.v.SetConfigFile(configFile)
	err := c.v.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		return err
	}

	return nil
}

// loadFromRemote 从远程kvstore中获取配置信息
func (c *ViperConfig) loadFromRemote() error {
	etcdUrl := c.v.GetString(c.option.Etcd.Url)
	if len(etcdUrl) == 0 {
		return errors.New("empty etcd address")
	}

	// 缺省的配置路径: <root>/<app>/<env>
	configPath := path.Join(c.option.Etcd.Root, c.app, c.env)
	err := c.v.AddRemoteProvider("etcd", etcdUrl, configPath)
	if err != nil {
		return errors.Wrapf(err, "add etcd provider, url: %s, path: %s", etcdUrl, configPath)
	}

	err = c.v.ReadRemoteConfig()
	if err != nil {
		return errors.Wrapf(err, "read etcd config, url: %s, path: %s", etcdUrl, configPath)
	}

	return nil
}

// getConfigFile 缺省的配置文件路径: <rootdir>/<basedir>/<app>/<app>.<env>.toml
func (c *ViperConfig) getConfigFile() string {
	var configFile string
	if c.env != "" {
		configFile = fmt.Sprintf("%s.%s.%s", c.app, c.env, c.option.File.Suffix)
	} else {
		configFile = fmt.Sprintf("%s.%s", c.app, c.option.File.Suffix)
	}
	return path.Join(c.option.File.RootDir, c.option.File.BaseDir, c.app, configFile)
}
