package ws

type ServerOption func(param *ServerParam)

type ServerParam struct {
	publicRouterGroup  *routerGroupParam
	protectRouterGroup *routerGroupParam
}

type routerGroupParam struct {
	Name      string
	UrlPrefix string
}

var (
	defaultServerParams = &ServerParam{
		publicRouterGroup: &routerGroupParam{
			Name:      "public",
			UrlPrefix: "/public",
		},
		protectRouterGroup: &routerGroupParam{
			Name:      "protect",
			UrlPrefix: "/api",
		},
	}
)

func WithPublicUrlPrefix(publicUrlPrefix string) ServerOption {
	return func(param *ServerParam) {
		param.publicRouterGroup.UrlPrefix = publicUrlPrefix
	}
}
