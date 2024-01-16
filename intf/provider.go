package intf

type Provider interface {
	Init(configer Configer, logger Logger, args ...any) error // 初始化
}
