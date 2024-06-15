package dapr

import (
	"fmt"
	reflectUtils "github.com/hdget/hdutils/reflect"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

type moduleKind int

const (
	moduleKindUnknown    moduleKind = iota
	ModuleKindInvocation            // dapr调用模块
	ModuleKindEvent                 // dapr事件模块
	ModuleKindDelayEvent            // 延迟事件模块
	ModuleKindHealth                // dapr健康检测模块
)

type moduler interface {
	GetApp() string
	GetModuleInfo() *moduleInfo
}

var (
	regModuleName        = regexp.MustCompile(`^[vV]([0-9]+)_([a-zA-Z0-9]+)`)
	errInvalidModule     = errors.New("invalid module, it must be struct")
	errInvalidModuleName = errors.New("invalid module name, it should be: v<number>_name, e,g: v1_abc")
	handlerNameSuffix    = "handler"
)

type baseModule struct {
	app        string      // 应用名称
	moduleInfo *moduleInfo // 模块的信息
}

type moduleInfo struct {
	StructName    string // 模块结构体的全名, 格式: "v<模块版本号>_<模块名>"
	ModuleName    string // 模块名
	ModuleVersion int    // 模块版本号
}

// newModule 从约定的结构名中解析模块名和版本, 结构名需要为v<number>_<module>
func newModule(app string, moduleObject any) (moduler, error) {
	structName := reflectUtils.GetStructName(moduleObject)
	if structName == "" {
		return nil, errInvalidModule
	}

	if !reflectUtils.IsAssignableStruct(moduleObject) {
		return nil, fmt.Errorf("module object: %s must be a pointer to struct", structName)
	}

	mInfo, err := parseModuleInfo(structName)
	if err != nil {
		return nil, err
	}

	return &baseModule{
		app:        app,
		moduleInfo: mInfo,
	}, nil
}

func (m *baseModule) GetApp() string {
	return m.app
}

// GetModuleInfo 获取模块元数据信息
func (m *baseModule) GetModuleInfo() *moduleInfo {
	return m.moduleInfo
}

func parseModuleInfo(structName string) (*moduleInfo, error) {
	tokens := regModuleName.FindStringSubmatch(structName)
	if len(tokens) != 3 {
		return nil, errInvalidModuleName
	}

	moduleVersion, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, errInvalidModuleName
	}

	return &moduleInfo{
		StructName:    structName,
		ModuleName:    tokens[2],
		ModuleVersion: moduleVersion,
	}, nil
}
