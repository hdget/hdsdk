package svc

import (
	"github.com/hdget/hdsdk/utils"
	"strconv"
)

// getModuleNameAndVersion 从约定的结构名中解析模块名和版本, 结构名需要为v<number>_<namespace>
func getModuleNameAndVersion(svcHolder any) (string, int, error) {
	moduleName := utils.GetStructName(svcHolder)
	tokens := regModuleName.FindStringSubmatch(moduleName)
	if len(tokens) != 3 {
		return "", 0, errInvalidModuleName
	}
	version, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", 0, errInvalidModuleName
	}
	return tokens[2], version, nil
}
