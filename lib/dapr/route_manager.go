package dapr

import (
	"bufio"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdsdk/v2/protobuf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RouteManager interface {
	DiscoverRouteItems(args ...HandlerNameMatcher) ([]*protobuf.RouteItem, error) // 获取路由项
	GetModulePath() string                                                        // 获取invocation module的路径
}

type RouteManagerImpl struct {
	baseDir       string
	relModulePath string
}

const (
	regexInvocationModulePath = `.+dapr.InvocationModule$`
)

// NewRouteManager 获取RouteManager
func NewRouteManager(baseDir string, skipDirs ...string) (RouteManager, error) {
	foundRelModulePath, err := findModulePath(baseDir, skipDirs...)
	if err != nil {
		return nil, err
	}

	if foundRelModulePath == "" {
		return nil, fmt.Errorf("no invocation module found in: %s", baseDir)
	}

	return &RouteManagerImpl{
		baseDir:       baseDir,
		relModulePath: foundRelModulePath,
	}, nil
}

func (rm RouteManagerImpl) GetModulePath() string {
	return rm.relModulePath
}

func (rm RouteManagerImpl) DiscoverRouteItems(handlerNameMatchers ...HandlerNameMatcher) ([]*protobuf.RouteItem, error) {
	routeItems := make([]*protobuf.RouteItem, 0)
	absModulePath := filepath.Join(rm.baseDir, rm.relModulePath)
	for _, moduleInstance := range LoadInvocationModules(rm.relModulePath) {
		routeAnnotations, err := moduleInstance.GetRouteAnnotations(absModulePath, handlerNameMatchers...)
		if err != nil {
			return nil, err
		}

		for _, ann := range routeAnnotations {
			for _, httpMethod := range ann.Methods {
				isPublic := int32(0)
				if ann.IsPublic {
					isPublic = 1
				}

				isRawResponse := int32(0)
				if ann.IsRawResponse {
					isRawResponse = 1
				}

				routeItems = append(routeItems, &protobuf.RouteItem{
					App:           moduleInstance.GetApp(),
					ModuleVersion: int32(moduleInstance.GetMeta().ModuleVersion),
					ModuleName:    moduleInstance.GetMeta().ModuleName,
					Handler:       ann.HandlerAlias,
					Endpoint:      ann.Endpoint,
					HttpMethod:    httpMethod,
					Permissions:   ann.Permissions,
					Origin:        ann.Origin,
					IsPublic:      isPublic,
					IsRawResponse: isRawResponse,
					Comment:       strings.Join(ann.Comments, "\r"),
				})
			}
		}
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
