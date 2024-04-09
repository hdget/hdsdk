package viper

import (
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"time"
)

type RemoteOption struct {
	provider          string
	url               string
	path              string
	disableRemoteEnvs []string // 禁用remote配置的环境列表
	watch             *RemoteWatchOption
}

type RemoteWatchOption struct {
	enabled     bool
	effectDelay int
}

var (
	defaultRemoteOption = &RemoteOption{
		provider:          "etcd3",                 // 默认的remote provider
		url:               "http://127.0.0.1:2379", // 默认的remote url
		disableRemoteEnvs: []string{"", "local"},   // 默认无环境或者local环境不需要加载remote配置,
		watch: &RemoteWatchOption{
			enabled:     true, // 默认是否检测远程配置变更
			effectDelay: 30,   // 配置变化生效时间为30秒
		},
	}
)

// LoadRemote 加载远程配置到变量configVar
func (vcLoader *viperConfigProvider) LoadRemote(configVar any) error {
	option := defaultRemoteOption

	// 当前环境不在disable列表时才需要加载remote配置
	if hdutils.Contains(option.disableRemoteEnvs, vcLoader.env) {
		return nil
	}

	// 尝试从远程配置信息
	err := vcLoader.loadFromRemote()
	if err != nil {
		hdutils.LogError("load config from remote", "err", err)
	}

	// 如果加载remote成功，则尝试监控配置变化
	if option.watch.enabled {
		err = vcLoader.watchRemote(configVar, option.watch)
		if err != nil {
			hdutils.LogError("watch remote config change", "err", err)
		}
	}

	return vcLoader.remote.Unmarshal(configVar)
}

//func (vcLoader *viperConfigProvider) getRemoteOption() *RemoteOption {
//	defaultRemoteOption.path = path.Join("/", path.Join(defaultValue.RootParts...), vcLoader.app) // 具体app的具体环境的配置保存在该路径下： /setting/app/<app>
//
//	// 尝试从sdk里面去取remoteOption配置
//	sdkConfiger, err := vcLoader.GetSDKConfig()
//	if err != nil {
//		return defaultRemoteOption
//	}
//
//	etcdConfig := sdkConfiger.GetEtcdConfig()
//	if etcdConfig == nil {
//		return defaultRemoteOption
//	}
//
//	if v := cast.ToString(etcdConfig["url"]); v != "" {
//		defaultRemoteOption.url = v
//	}
//
//	if v := cast.ToString(etcdConfig["path"]); v != "" {
//		defaultRemoteOption.url = v
//	}
//	return defaultRemoteOption
//}

// loadFromRemote 尝试从远程kvstore中获取配置信息
// windows下测试: e,g: type test.txt | etcdctl.exe put /setting/app/hello/test
func (vcLoader *viperConfigProvider) loadFromRemote() error {
	if len(vcLoader.remoteOptions) == 0 {
		//vcLoader.remoteOptions = append(vcLoader.remoteOptions, vcLoader.getRemoteOption())
		vcLoader.remoteOptions = append(vcLoader.remoteOptions, defaultRemoteOption)
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

func (vcLoader *viperConfigProvider) watchRemote(remoteConfigVar any, option *RemoteWatchOption) error {
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
			time.Sleep(time.Second * time.Duration(option.effectDelay)) // delay after each request

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

//// UpdateRemoteConfig 更新远程配置
//// nolint: staticcheck
//func (vcLoader *viperConfigProvider) UpdateRemoteConfig(v any) error {
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
