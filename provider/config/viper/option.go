package viper

import (
	"github.com/hdget/hdutils"
	"github.com/spf13/viper"
	"path"
)

type Option func(loader *viperConfigProvider)

type fileOption struct {
	configFile string   // 指定的配置文件
	dirs       []string // 如果未指定配置文件情况下，搜索的目录
	filename   string   // 如果未指定配置文件情况下，搜索的文件名，不需要文件后缀
}

func WithConfigFile(filepath string) Option {
	return func(c *viperConfigProvider) {
		c.fileOption.configFile = filepath
	}
}

func WithEnvPrefix(envPrefix string) Option {
	return func(c *viperConfigProvider) {
		c.envPrefix = envPrefix
	}
}

func WithConfigDir(args ...string) Option {
	return func(c *viperConfigProvider) {
		c.fileOption.dirs = append(c.fileOption.dirs, args...)
	}
}

func WithConfigFilename(filename string) Option {
	return func(c *viperConfigProvider) {
		if path.Ext(filename) != "" {
			hdutils.LogWarn("filename should not contains suffix", "filename", filename)
		}
		c.fileOption.filename = filename
	}
}

func WithConfigType(configType string) Option {
	return func(c *viperConfigProvider) {
		if !hdutils.Contains(viper.SupportedExts, configType) {
			hdutils.LogFatal("set configer type", "supported", viper.SupportedExts, "err", viper.UnsupportedConfigError(configType))
		}
		c.configType = configType
	}
}

func WithRemote(provider, url, path string) Option {
	return func(c *viperConfigProvider) {
		if !hdutils.Contains(viper.SupportedRemoteProviders, provider) {
			hdutils.LogFatal("set remote configer provider", "supported", viper.SupportedRemoteProviders, "err", viper.UnsupportedRemoteProviderError(provider))
		}

		if c.remoteOptions == nil {
			c.remoteOptions = make([]*RemoteOption, 0)
		}

		c.remoteOptions = append(c.remoteOptions, &RemoteOption{
			provider: provider,
			url:      url,
			path:     path,
		})
	}
}

func WithRoot(args ...string) Option {
	return func(c *viperConfigProvider) {
		rootParts := make([]string, 0)
		c.rootParts = append(rootParts, args...)
	}
}

func WithConfigContent(content string) Option {
	return func(c *viperConfigProvider) {
		c.content = hdutils.StringToBytes(content)
	}
}
