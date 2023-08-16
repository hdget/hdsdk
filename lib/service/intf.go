package service

type Module interface {
	GetName() string
	Super(m any)
	GetRoutes(srcPath string) ([]*Route, error)
	GetHandlers(args ...MethodMatchFunction) map[string]any
	//GetPermGroups(srcPath string) ([]*PermGroup, error)
}

type Option func(module Module) Module

type MethodMatchFunction func(methodName string) (newMethodName string, matched bool) // 传入receiver.methodName, 判断是否匹配，然后取出处理后的method名
