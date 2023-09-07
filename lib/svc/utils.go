package svc

import (
	"github.com/hdget/hdsdk/utils"
	"regexp"
	"strconv"
)

var (
	regModuleName = regexp.MustCompile(`^[vV]([0-9]+)_([a-zA-Z0-9]+)`)
)

// getModuleInfo 从约定的结构名中解析模块名和版本, 结构名需要为v<number>_<module>
func getModuleInfo(svcHolderOrModuleName any) (*moduleInfo, error) {
	var structName string
	switch v := svcHolderOrModuleName.(type) {
	case string:
		structName = v
	default:
		structName = utils.Reflect().GetStructName(v)
	}

	tokens := regModuleName.FindStringSubmatch(structName)
	if len(tokens) != 3 {
		return nil, errInvalidModuleName
	}
	version, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, errInvalidModuleName
	}

	return &moduleInfo{
		Name:    structName,
		Module:  tokens[2],
		Version: version,
	}, nil

}
