package hddapr

import (
	"github.com/hdget/hdutils"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

var (
	regModuleName    = regexp.MustCompile(`^[vV]([0-9]+)_([a-zA-Z0-9]+)`)
	errInvalidModule = errors.New("invalid module, it must be struct")
)

type ModuleInfo struct {
	Name          string // 结构名, 格式: "v<模块版本号>_<模块名>"
	ModuleName    string // 模块名
	ModuleVersion int    // 模块版本号
}

// parseModuleInfo 从约定的结构名中解析模块名和版本, 结构名需要为v<number>_<module>
func parseModuleInfo(moduleObject any) (*ModuleInfo, error) {
	structName := hdutils.Reflect().GetStructName(moduleObject)
	if structName == "" {
		return nil, errInvalidModule
	}
	return toModuleInfo(structName)
}

// toModuleInfo 将结构名转换为模块信息
func toModuleInfo(structName string) (*ModuleInfo, error) {
	tokens := regModuleName.FindStringSubmatch(structName)
	if len(tokens) != 3 {
		return nil, errInvalidModuleName
	}
	version, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, errInvalidModuleName
	}

	return &ModuleInfo{
		Name:          structName,
		ModuleName:    tokens[2],
		ModuleVersion: version,
	}, nil
}
