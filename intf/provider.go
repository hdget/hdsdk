package intf

// Provider 底层库能力提供者接口
type Provider interface {
	Init(logger Logger, args ...any) error // 初始化
}
