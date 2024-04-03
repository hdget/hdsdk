package hddapr

//
////// genHandlerId 生成handlerId
////func (b *invocationModuleImpl) genHandlerId(moduleName, methodName string) string {
////	return strings.Join([]string{moduleName, methodName}, "_")
////}
//
//// GetHandlers 将map[string]*invocationHandler转换成map[string]common.serviceHandler
//func (m *invocationModuleImpl) GetHandlers() map[string]any {
//	handlers := make(map[string]any)
//	// h为*invocationHandler
//	for _, h := range m.functions {
//		// daprMethodName = v2:xxx:handlerName
//		daprMethodName := GetServiceInvocationName(m.ModuleVersion, m.Module, h.name)
//		daprMethod := m.toDaprServiceInvocationHandler(h.handler)
//		if daprMethod != nil {
//			handlers[daprMethodName] = daprMethod
//		}
//
//	}
//	return handlers
//}
//
//func (m *invocationModuleImpl) ValidateHandler(handler any) error {
//	if hdutils.Reflect().GetFuncSignature(handler) != hdutils.Reflect().GetFuncSignature(InvocationFunction(nil)) {
//		return fmt.Errorf("invalid handler: %s, it should be: func(ctx context.Context, event *common.InvocationEvent) (any, error)", hdutils.Reflect().GetFuncName(handler))
//	}
//	return nil
//}
