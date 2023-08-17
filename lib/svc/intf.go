package svc

type Module interface {
	GetName() string
	GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error)
	DiscoverHandlers(args ...HandlerMatch) (map[string]any, error) // 通过反射发现Handlers
	RegisterHandlers(handlers map[string]any)
	GetHandlers() map[string]any // 获取手动注册的handlers
	//GetPermGroups(srcPath string) ([]*PermGroup, error)
}
