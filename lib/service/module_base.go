package service

import "github.com/dapr/go-sdk/service/common"

type BaseModule struct {
	App       string
	Version   int                                        // 版本号
	Namespace string                                     // 命名空间
	Name      string                                     // 服务模块名
	Handlers  map[string]common.ServiceInvocationHandler // 定义的方法

	DeferFunc        func(string) // defer functions before handler
	EnterFunctions   []Func
	SuccessFunctions []Func
	FailedFunctions  []Func
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

func (b *BaseModule) GetApp() string {
	return b.App
}
