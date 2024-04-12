package intf

type ConfigLoader interface {
	Unmarshal(configVar any, key ...string) error
}
