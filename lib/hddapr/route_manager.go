package hddapr

import (
	"bufio"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdsdk/lib/ws"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RouteManager interface {
	GetRouteItems(args ...HandlerNameMatcher) ([]*RouteItem, error) // 获取路由项
	GetModulePath() string                                          // 获取模块路径
}

type RouteManagerImpl struct {
	baseDir       string
	relModulePath string
}

type RouteItem struct {
	App           string        // app name
	ModuleName    string        // module name
	ModuleVersion int           // module version
	Handler       string        // handler alias
	Endpoint      string        // endpoint
	HttpMethod    ws.HttpMethod // http methods
	Origin        string        // 请求来源
	IsRawResponse bool          // 是否返回原始消息
	IsPublic      bool          // 是否是公共路由
	Url           string        // 完整路由
	Permissions   []string      // 所属权限列表
	Comments      []string      // 备注
}

const (
	regexInvocationModulePath = `.+hddapr.InvocationModule$`
)

// NewRouteManager 获取RouteManager
func NewRouteManager(baseDir string, skipDirs ...string) (RouteManager, error) {
	relModulePath, err := findModulePath(baseDir, skipDirs...)
	if err != nil {
		return nil, err
	}

	if relModulePath == "" {
		return nil, fmt.Errorf("no invocation module found in: %s", baseDir)
	}

	return &RouteManagerImpl{
		baseDir:       baseDir,
		relModulePath: relModulePath,
	}, nil
}

func (rm RouteManagerImpl) GetModulePath() string {
	return rm.relModulePath
}

func (rm RouteManagerImpl) GetRouteItems(handlerNameMatchers ...HandlerNameMatcher) ([]*RouteItem, error) {
	routeItems := make([]*RouteItem, 0)
	absModulePath := filepath.Join(rm.baseDir, rm.relModulePath)
	for moduleName, moduleInstance := range LoadInvocationModules(rm.relModulePath) {
		routeAnnotations, err := moduleInstance.GetRouteAnnotations(absModulePath, handlerNameMatchers...)
		if err != nil {
			return nil, err
		}

		handlerNames := make([]string, 0)
		for _, ann := range routeAnnotations {
			for _, httpMethod := range ann.Methods {
				routeItems = append(routeItems, &RouteItem{
					App:           moduleInstance.GetApp(),
					ModuleVersion: moduleInstance.GetInfo().ModuleVersion,
					ModuleName:    moduleInstance.GetInfo().ModuleName,
					Handler:       ann.HandlerAlias,
					Endpoint:      ann.Endpoint,
					HttpMethod:    ws.ToHttpMethod(httpMethod),
					Permissions:   ann.Permissions,
					Origin:        ann.Origin,
					IsPublic:      ann.IsPublic,
					IsRawResponse: ann.IsRawResponse,
					Comments:      ann.Comments,
				})
			}

			handlerNames = append(handlerNames, ann.HandlerAlias)
		}

		fmt.Printf(" - module: %-25s total: %-5d functions: [%s]\n", moduleName, len(routeAnnotations), strings.Join(handlerNames, ", "))
	}

	return routeItems, nil
}

// findModulePath 尝试找到
func findModulePath(baseDir string, skipDirs ...string) (string, error) {
	st, err := os.Stat(baseDir)
	if err != nil {
		return "", err
	}

	if !st.IsDir() {
		return "", fmt.Errorf("invalid dir, dir: %s", baseDir)
	}

	var found string
	match, _ := regexp.Compile(regexInvocationModulePath)
	_ = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && pie.Contains(skipDirs, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				d := scanner.Text()
				if match.MatchString(d) {
					parentDir, _ := filepath.Split(path)
					relDir, _ := filepath.Rel(baseDir, parentDir)
					found = relDir
					break
				}
			}
		}
		return nil
	})
	return found, nil
}
