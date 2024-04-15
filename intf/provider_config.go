package intf

type ConfigProvider interface {
	Unmarshal(configVar any, key ...string) error
}
