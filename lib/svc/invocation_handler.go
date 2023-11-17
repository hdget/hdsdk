package svc

import (
	"github.com/hdget/hdutils"
	"strings"
)

type invocationHandler struct {
	id     string //
	method string // receiver.method对应的方法名，例如： (*aaa) GetIdHandler(), 这里GetIdHandler为method名
	alias  string // 别名
	fn     any    // 具体的调用函数
}

func newInvocationHandler(moduleName, alias string, fn any) *invocationHandler {
	methodName := hdutils.Reflect().GetFuncName(fn)
	return &invocationHandler{
		id:     genHandlerId(moduleName, methodName),
		method: methodName,
		alias:  alias,
		fn:     fn,
	}
}

// genHandlerId 生成handlerId
func genHandlerId(moduleName, methodName string) string {
	return strings.Join([]string{moduleName, methodName}, "_")
}
