package svc

type Module interface {
	GetName() string
	Super(m any)
	GetRoutes(srcPath string, args ...HandlerMatch) ([]*Route, error)
	GetServiceHandlers(args ...HandlerMatch) (map[string]any, error)
	//GetPermGroups(srcPath string) ([]*PermGroup, error)
}

type Option func(module Module) Module
