package types

// 环境变量的参数
const (
	ENV_PRODUCTION     = "prod"
	ENV_PRE_PRODUCTION = "pre"
	ENV_SIMULATION     = "sim"
	ENV_TEST           = "test"
	ENV_DEV            = "dev"
	ENV_LOCAL          = "local"
)

// 所有支持的环境
var SupportedEnvs = []string{
	ENV_PRODUCTION,
	ENV_PRE_PRODUCTION,
	ENV_SIMULATION,
	ENV_TEST,
	ENV_DEV,
	ENV_LOCAL,
}
