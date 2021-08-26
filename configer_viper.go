package hdsdk

import (
	"fmt"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"path"
)

// 配置加载选项
type ConfigOption struct {
	Env  EnvOption  // 环境变量选项
	File FileOption // 配置文件选项
	Etcd EtcdOption // etcd选项
}

type EnvOption struct {
	Prefix string // 环境变量前缀
	Name   string // 环境变量名字
	Value  string // 环境变量值
}

type FileOption struct {
	RootDir string // 配置文件所在的根目录
	BaseDir string // 配置文件的上级目录
	Suffix  string // .toml
}

type EtcdOption struct {
	Url string
}

// 缺省配置选项
var (
	defaultEnvOption = EnvOption{
		Prefix: "HDGET",
		Name:   "RUNTIME",
		Value:  types.ENV_PRODUCTION,
	}

	defaultFileOption = FileOption{
		RootDir: "setting",
		BaseDir: "app",
		Suffix:  "toml",
	}

	defaultEtcdOption = EtcdOption{
		Url: "http://127.0.0.1:2379",
	}

	defaultConfigOption = ConfigOption{
		Env:  defaultEnvOption,
		File: defaultFileOption,
		Etcd: defaultEtcdOption,
	}
)

// LoadConfig It will will load config from difference sources,
// the config will read with following precedence:
// - flag(命令行的设置是优先级最高的)
// - env(其次为环境变量)
// - config(再其次为配置文件)
// - key/value store(接着是KVStore(etcd)里面的设置)
// - default(最低的是默认配置)
func LoadConfig(app, cliEnv, cliFile string, args ...ConfigOption) *viper.Viper {
	// 获取加载选项
	option := defaultConfigOption
	if len(args) > 0 {
		option = args[0]
	}

	v := viper.New()

	// 默认配置
	loadFromDefault(v)

	// 尝试从环境中读取配置
	loadFromEnv(v, option)

	// 设置有效的环境变量
	envValue, err := setupEnv(v, cliEnv, option)
	if err != nil {
		utils.Fatal("setup env value", "err", err)
	}

	// 尝试从远程KV store，例如etcd加载配置信息
	err = loadFromRemote(v, app, envValue, option)
	if err != nil {
		utils.Print("ERR", "load config from etcd", "err", err)
	}

	// 尝试从配置中读取配置信息
	absPath, err := loadFromFile(v, app, envValue, cliFile)
	if err != nil {
		utils.Print("ERR", "load config from file", "file", absPath, "err", err)
	}

	return v
}

/////////////////////////////////////////////////////////////////
// private functions
////////////////////////////////////////////////////////////////
func loadFromDefault(v *viper.Viper) {
	v.SetDefault("env", defaultEnvOption.Value)
	v.SetDefault("config_url", defaultEtcdOption.Url)
}

// loadFromEnv 从环境文件中读取配置信息
func loadFromEnv(v *viper.Viper, option ConfigOption) {
	// 获取有效的环境变量prefix
	envPrefix := option.Env.Prefix
	if envPrefix == "" {
		envPrefix = defaultEnvOption.Prefix
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
}

// 设置环境变量, 首先从命令行读取env, 如果未读到合法的env，尝试从环境变量runtime中去获取环境变量
func setupEnv(v *viper.Viper, cliEnv string, option ConfigOption) (string, error) {
	envValue := cliEnv
	if exist := utils.StringSliceContains(types.SupportedEnvs, envValue); !exist {
		// 从环境变量通过读取env_name来获取env的值
		envValue = v.GetString(option.Env.Name)
	}

	if exist := utils.StringSliceContains(types.SupportedEnvs, envValue); !exist {
		return "", errors.Errorf("invalid env: %s, supported envs is: %v", envValue, types.SupportedEnvs)
	}

	v.Set("ENV", envValue)
	return envValue, nil
}

// getDefaultConfigFile 缺省的配置文件路径: <rootdir>/setting/app/<app>.<env>.toml
func getDefaultConfigFile(app, envValue string) string {
	configFile := fmt.Sprintf("%s.%s.%s", app, envValue, defaultFileOption.Suffix)
	return path.Join(defaultFileOption.RootDir, defaultFileOption.BaseDir, configFile)
}

// 从配置文件中读取配置信息
func loadFromFile(v *viper.Viper, app, envValue, cliFile string) (string, error) {
	configFile := cliFile
	if configFile == "" {
		configFile = getDefaultConfigFile(app, envValue)
	}

	// optionally look for config in the working directory
	v.AddConfigPath(".")
	v.SetConfigFile(configFile)
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		return configFile, err
	}

	return configFile, nil
}

// loadFromRemote 从远程kvstore中获取配置信息
func loadFromRemote(v *viper.Viper, app, env string, option ConfigOption) error {
	etcdUrl := v.GetString(option.Etcd.Url)
	if len(etcdUrl) == 0 {
		return errors.New("empty etcd address")
	}

	// 缺省的配置路径: <rootdir>/<app>/<env>
	configPath := path.Join(option.File.RootDir, app, env)
	err := v.AddRemoteProvider("etcd", etcdUrl, configPath)
	if err != nil {
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		return err
	}

	return nil
}
