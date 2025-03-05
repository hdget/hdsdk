package dapr

import "os"

const (
	_envVarNamespace = "HD_NAMESPACE"
)

// injectEnv 根据环境变量调整输入字符串。
// 该函数检查名为 _envVarNamespace 的环境变量是否存在。
// 如果存在，将环境变量的值与输入字符串用下划线连接后返回；
// 如果不存在，直接返回输入字符串。
func injectEnv(input string) string {
	if namespace, exists := os.LookupEnv(_envVarNamespace); exists {
		return namespace + "_" + input
	}
	return input
}
