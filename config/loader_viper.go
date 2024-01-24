package config

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/fx"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type vipConfigLoaderParams struct {
	fx.In
	app     string
	env     string
	options []Option
}

// viperConfigLoader 命令行配置
type viperConfigLoader struct {
	mu                sync.RWMutex
	local             *viper.Viper
	remote            *viper.Viper
	app               string
	env               string
	envPrefix         string          // 环境变量前缀
	rootParts         []string        // 根目录的部分，例如/setting/app, 则为[]string{"configer", "app"}
	configType        string          // 配置内容类型，e,g: toml, json
	disableRemoteEnvs []string        // 禁用remote配置的环境列表
	fileOption        *fileOption     // 文件配置选项
	remoteOptions     []*remoteOption // 远程加载配置选项
	watchOption       *watchOption    // 检测配置变化选项
	content           []byte          // 如果用WithConfigContent指定了配置内容，则这里不为空
}

// 缺省配置选项
var (
	defaultValue = struct {
		EnvPrefix         string
		RootParts         []string
		ConfigType        string
		FileOption        *fileOption
		RemoteUrlKey      string
		RemoteProvider    string
		RemoteUrl         string
		WatchOption       *watchOption
		DisableRemoteEnvs []string
	}{
		EnvPrefix:  "HD",
		RootParts:  []string{"setting", "app"}, // 缺省的Root路径Parts, 这里定义为[]string，方便用path.Join或者filepath.Join
		ConfigType: "toml",                     // 缺省的配置文件类型
		FileOption: &fileOption{
			dirs: make([]string, 0),
		},
		RemoteUrlKey:   "sdk.etcd.url",          // 默认etcd在文件中定义的key
		RemoteProvider: "etcd3",                 // 默认的remote provider
		RemoteUrl:      "http://127.0.0.1:2379", // 默认的remote url
		WatchOption: &watchOption{
			enabled:     true, // 默认是否检测远程配置变更
			effectDelay: 30,   // 配置变化生效时间为30秒
		},
		DisableRemoteEnvs: []string{"", "local"}, // 默认无环境或者local环境不需要加载remote配置
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

// NewConfigLoader 初始化config loader
func NewConfigLoader(app, env string, options ...Option) intf.ConfigLoader {
	loader := &viperConfigLoader{
		local:             viper.New(),
		remote:            viper.New(),
		app:               app,
		env:               env,
		envPrefix:         defaultValue.EnvPrefix,
		rootParts:         defaultValue.RootParts,
		configType:        defaultValue.ConfigType,
		disableRemoteEnvs: defaultValue.DisableRemoteEnvs, // 禁用remote配置加载的环境列表
		fileOption:        defaultValue.FileOption,
		remoteOptions:     nil,
		watchOption:       defaultValue.WatchOption,
	}

	for _, option := range options {
		option(loader)
	}

	return loader
}

// LoadKey 将key对应的配置unmarshal到变量configVar
func (vcLoader *viperConfigLoader) LoadKey(key string, configVar any) error {
	err := vcLoader.loadLocal()
	if err != nil {
		return err
	}

	return vcLoader.local.UnmarshalKey(key, &configVar)
}

// Load 将除key以外的配置unmarshal到变量configVar
func (vcLoader *viperConfigLoader) Load(configVar any) error {
	err := vcLoader.loadLocal()
	if err != nil {
		return err
	}
	return vcLoader.local.Unmarshal(&configVar)
}

//// UpdateRemoteConfig 更新远程配置
//// nolint: staticcheck
//func UpdateRemoteConfig(v any) error {
//	if hdsdk.Etcd == nil {
//		return errors.New("hdsdk not initialized")
//	}
//
//	if _vc == nil {
//		return errors.New("configer not initialized")
//	}
//
//	data, err := json.Marshal(v)
//	if err != nil {
//		return err
//	}
//
//	for _, option := range _vc.remoteOptions {
//		err = hdsdk.Etcd.Set(option.path, data)
//		if err != nil {
//			return err
//		}
//		// 如果成功
//		break
//	}
//
//	return nil
//}

//func (vcLoader *viperConfigLoader) Read(data []byte) *viperConfigLoader {
//	vcLoader.local.SetConfigType(vcLoader.configType)
//	_ = vcLoader.local.MergeConfig(bytes.NewReader(data))
//	return vcLoader
//}
//
//func (vcLoader *viperConfigLoader) ReadString(s string) *viperConfigLoader {
//	_ = vcLoader.Read(hdutils.StringToBytes(s))
//	return vcLoader
//}

// ///////////////////////////////////////////////////////////////
// private functions
// //////////////////////////////////////////////////////////////

// Load 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中）
// - env: 环境变量配置(高)
func (vcLoader *viperConfigLoader) loadLocal() error {
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

func (vcLoader *viperConfigLoader) loadRemote() error {
	// 当前环境不在disable列表时才需要加载remote配置
	if hdutils.Contains(vcLoader.disableRemoteEnvs, vcLoader.env) {
		return nil
	}

	// 尝试从远程配置信息
	err := vcLoader.loadFromRemote()
	if err != nil {
		hdutils.LogError("load configloader from remote", "err", err)
	}

	//// 如果加载remote成功，则尝试监控配置变化
	//if vcLoader.watchOption.enabled {
	//	err = vcLoader.watchRemote(hdsdk.configer)
	//	if err != nil {
	//		hdutils.LogError("watch remote configloader change", "err", err)
	//	}
	//}
	//
	return nil
}

func (vcLoader *viperConfigLoader) watchRemote(remoteConfigVar any) error {
	// 如果无任何远程配置设置，忽略
	if len(vcLoader.remoteOptions) == 0 {
		return nil
	}

	// currently, only tested with etcd support
	err := vcLoader.remote.WatchRemoteConfigOnChannel()
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(vcLoader.watchOption.effectDelay)) // delay after each request

			// 加写锁保证remoteConfigVar没有同时被写
			vcLoader.mu.Lock()
			err = vcLoader.remote.Unmarshal(remoteConfigVar)
			vcLoader.mu.Unlock()
			if err != nil {
				hdutils.LogError("unable to unmarshal remote configer", "err", err)
			}
		}
	}()
	return nil
}

// loadFromEnv 从环境文件中读取配置信息
func (vcLoader *viperConfigLoader) loadFromEnv() {
	// 如果设置了环境变量前缀，则尝试自动获取环境变量中的配置
	if vcLoader.envPrefix != "" {
		vcLoader.local.SetEnvPrefix(vcLoader.envPrefix)
		vcLoader.local.AutomaticEnv()
	}
}

func (vcLoader *viperConfigLoader) loadMinimal() error {
	minimalConfig := fmt.Sprintf(tplMinimalConfigContent, vcLoader.app)
	return vcLoader.local.MergeConfig(bytes.NewReader(hdutils.StringToBytes(minimalConfig)))
}

func (vcLoader *viperConfigLoader) loadFromFile() error {
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

// loadFromRemote 尝试从远程kvstore中获取配置信息
// windows下测试: e,g: type test.txt | etcdctl.exe put /setting/app/hello/test
func (vcLoader *viperConfigLoader) loadFromRemote() error {
	if len(vcLoader.remoteOptions) == 0 {
		vcLoader.remoteOptions = append(vcLoader.remoteOptions, vcLoader.getDefaultRemoteOption())
	}

	for _, option := range vcLoader.remoteOptions {
		err := vcLoader.remote.AddRemoteProvider(option.provider, option.url, option.path)
		if err != nil {
			return errors.Wrapf(err, "add remote provider, provider: %s, url: %s, path: %s", option.provider, option.url, option.path)
		}
	}

	// 远程的固定为json
	vcLoader.remote.SetConfigType("json")
	err := vcLoader.remote.ReadRemoteConfig()
	if err != nil {
		return errors.Wrapf(err, "read remote configer")
	}

	for _, option := range vcLoader.remoteOptions {
		hdutils.LogDebug("load configer from remote", "provider", option.provider, "url", option.url, "path", option.path)
	}
	return nil
}

func (vcLoader *viperConfigLoader) getDefaultRemoteOption() *remoteOption {
	// 加载远程配置的时候优先从之前已经读取的配置，例如文件配置中取remoteUrl
	url := defaultValue.RemoteUrl
	if v := vcLoader.local.GetString(defaultValue.RemoteUrlKey); v != "" {
		url = v
	}

	return &remoteOption{
		provider: defaultValue.RemoteProvider,
		url:      url,
		path:     path.Join("/", path.Join(defaultValue.RootParts...), vcLoader.app), // 具体app的具体环境的配置保存在该路径下： /setting/app/<app>

	}
}
