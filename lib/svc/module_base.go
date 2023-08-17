package svc

type baseModule struct {
	self     any // 实际module
	App      string
	Version  int    // 版本号
	Name     string // 服务模块名
	handlers map[string]any
}

func newBaseModule(app string, svcHolder any) (*baseModule, error) {
	moduleName, version, err := getModuleNameAndVersion(svcHolder)
	if err != nil {
		return nil, err
	}

	return &baseModule{
		App:     app,
		Version: version,
		Name:    moduleName,
		self:    svcHolder,
	}, nil
}

func (b *baseModule) GetName() string {
	return b.Name
}

func (b *baseModule) RegisterHandlers(handlers map[string]any) {
	b.handlers = handlers
	addRegistry(b.Name, b.self.(Module))
}
