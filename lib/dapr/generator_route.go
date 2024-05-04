package dapr

//
//// route generator
//type routeGeneratorImpl struct {
//	*baseGeneratorImpl
//	app               string
//	srcPath           string
//	invocationModules map[string]InvocationModule
//}
//
//type RouteItem struct {
//	App           string
//	ModuleVersion int
//	ModuleName        string
//	HandlerAlias       string
//	Endpoint      string
//	Origin        string
//	IsPublic      bool
//	IsRawResponse bool
//	HttpMethods   []string
//	Permissions   []string
//	Comments      []string
//}
//
//func NewRouteGenerator(app, srcPath string, invocationModules map[string]InvocationModule) Generator {
//	return &routeGeneratorImpl{
//		baseGeneratorImpl: &baseGeneratorImpl{},
//		app:               app,
//		srcPath:           srcPath,
//		invocationModules: invocationModules,
//	}
//}
//
//func (m *routeGeneratorImpl) Gen() error {
//	// 获取routeItems
//	routeItems := make([]any, 0)
//	fmt.Printf("Generating routes from: %s...\n", m.srcPath)
//	for moduleName, module := range m.invocationModules {
//		routeAnnotations, err := module.GetRouteAnnotations(m.srcPath)
//		if err != nil {
//			return err
//		}
//
//		handlerNames := make([]string, 0)
//		for _, ann := range routeAnnotations {
//			routeItems = append(routeItems, &RouteItem{
//				App:           m.app,
//				ModuleVersion: ann.ModuleVersion,
//				ModuleName:        ann.ModuleName,
//				HandlerAlias:       ann.HandlerAlias,
//				Endpoint:      ann.Endpoint,
//				HttpMethods:   ann.HttpMethods,
//				Permissions:   ann.Permissions,
//				Origin:        ann.Origin,
//				IsPublic:      ann.IsPublic,
//				IsRawResponse: ann.IsRawResponse,
//				Comments:      ann.Comments,
//			})
//
//			handlerNames = append(handlerNames, ann.HandlerAlias)
//		}
//
//		fmt.Printf(" - module: %-25s total: %-5d functions: [%s]\n", moduleName, len(routeAnnotations), strings.Join(handlerNames, ", "))
//	}
//
//	// 获取当前函数所在的包
//	pc, _, _, _ := runtime.Caller(0)
//	splitFuncName := strings.Split(runtime.FuncForPC(pc).ExchangeName(), ".")
//	packagePath := strings.Join(splitFuncName[0:len(splitFuncName)-2], ".")
//
//	varName := "routes"
//	return hdutils.
//		NewGoFile("autogen", map[string]string{packagePath: ""}).
//		DeclareSliceVar(varName, packagePath, routeItems).
//		AddMethod(hdutils.Reflect().GetStructName(m), hdutils.Reflect().GetFuncName(m.Get), nil, []string{"any"}, []jen.Code{jen.Return(jen.Id(varName))}).
//		Save(path.Join("autogen", "routes.go"))
//}

//
//func (m *routeGeneratorImpl) Register() error {
//	generated := m.Get()
//	if generated == nil {
//		return errors.New("invalid generated stuff")
//	}
//
//	routeItems := generated.([]*RouteItem)
//	if len(routeItems) == 0 {
//		return nil
//	}
//
//	data, _ := json.Marshal(routeItems)
//	_, err := Invoke("gateway", 2, "route", "update", data)
//	return err
//}
