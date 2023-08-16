package svc

type BaseModule struct {
	realModule any // 实际module
	App        string
	Version    int    // 版本号
	Namespace  string // 命名空间
	Name       string // 服务模块名
}

func NewBaseModule(app, name string, version int) *BaseModule {
	return &BaseModule{
		App:     app,
		Version: version,
		Name:    name,
	}
}

func (b *BaseModule) GetName() string {
	return b.Name
}

func (b *BaseModule) Super(m any) {
	b.realModule = m
}
