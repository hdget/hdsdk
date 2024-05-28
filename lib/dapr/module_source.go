package dapr

import (
	"bytes"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type ModuleSourceHandler interface {
	Handle() (*ModuleSource, error)
}

// ModuleSource 模块源代码信息
type ModuleSource struct {
	Routes           []*protobuf.RouteItem
	ModulePaths      map[ModuleKind]string // 模块的路径
	ServerSourceFile string                // dapr.NewGrpcServer所在的go文件
}

type moduleAstParserImpl struct {
	sourceRootDir       string
	skipDirs            []string
	handlerNameMatchers []HandlerNameMatcher
}

type ModuleSourceHandleOption func(*moduleAstParserImpl)

const (
	exprInvocationModule = "&{dapr InvocationModule}"
	exprEventModule      = "&{dapr EventModule}"
	libDaprImportPath    = "github.com/hdget/hdsdk/v2/lib/dapr"
)

var (
	daprServerFunctions = []string{
		"NewGrpcServer",
		"NewHttpServer",
	}
)

// NewModuleSourceHandler 获取模块源代码处理器
func NewModuleSourceHandler(baseDir string, options ...ModuleSourceHandleOption) ModuleSourceHandler {
	p := &moduleAstParserImpl{
		sourceRootDir: baseDir,
	}
	for _, option := range options {
		option(p)
	}

	return p
}

func WithSkipDirs(skipDirs ...string) ModuleSourceHandleOption {
	return func(impl *moduleAstParserImpl) {
		impl.skipDirs = skipDirs
	}
}

func WithHandlerMatchers(handlerNameMatchers ...HandlerNameMatcher) ModuleSourceHandleOption {
	return func(impl *moduleAstParserImpl) {
		impl.handlerNameMatchers = handlerNameMatchers
	}
}

func (p moduleAstParserImpl) Handle() (*ModuleSource, error) {
	moduleSource, err := p.parseModuleSourceInfo(p.skipDirs...)
	if err != nil {
		return nil, err
	}

	// 如果找到dapr.NewGrpcServer或者dapr.NewHttpServer则需要将导入invocationModule和eventModule
	if moduleSource.ServerSourceFile != "" && len(moduleSource.ModulePaths) > 0 {
		err = p.addImportModulePaths(moduleSource.ServerSourceFile, moduleSource.ModulePaths)
		if err != nil {
			return nil, err
		}
	}

	// 如果找到invocationModule则需要查找路由
	if moduleSource.ModulePaths[ModuleKindInvocation] != "" {
		moduleSource.Routes, err = p.discoverRouteItems(moduleSource.ModulePaths[ModuleKindInvocation])
		if err != nil {
			return nil, err
		}
	}

	return moduleSource, nil
}

// parseModuleSourceInfo parse source info
func (p moduleAstParserImpl) parseModuleSourceInfo(skipDirs ...string) (*ModuleSource, error) {
	st, err := os.Stat(p.sourceRootDir)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		return nil, fmt.Errorf("invalid source code dir, dir: %s", p.sourceRootDir)
	}

	sourceInfo := &ModuleSource{
		ModulePaths: make(map[ModuleKind]string),
	}
	_ = filepath.Walk(p.sourceRootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// 忽略掉.开头的隐藏目录和需要skip的目录
			if strings.HasPrefix(info.Name(), ".") || pie.Contains(skipDirs, info.Name()) {
				return filepath.SkipDir
			}
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return errors.Wrapf(err, "golang ast parse file, path: %s", path)
			}

			// 新建一个记录导入别名与包名映射关系的字典
			importAliase2importPath := make(map[string]string)
			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.GenDecl:
					// 仅处理类型声明
					if n.Tok == token.TYPE {
						for _, spec := range n.Specs {
							// 如果类型规范是类型别名或类型声明
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								// 如果类型规范是结构体类型
								structures, ok := typeSpec.Type.(*ast.StructType)
								if ok {
									// 检查第一个field是否是匿名引入的模块， e,g: type A struct { dapr.InvocationModule }
									if len(structures.Fields.List) > 0 {
										switch fmt.Sprintf("%s", structures.Fields.List[0].Type) {
										case exprInvocationModule:
											found, _ := filepath.Rel(p.sourceRootDir, filepath.Dir(path))
											sourceInfo.ModulePaths[ModuleKindInvocation] = filepath.ToSlash(found)
										case exprEventModule:
											found, _ := filepath.Rel(p.sourceRootDir, filepath.Dir(path))
											sourceInfo.ModulePaths[ModuleKindEvent] = filepath.ToSlash(found)
										}
									}
								}

							}
						}
					}
				case *ast.SelectorExpr:
					if x, ok := n.X.(*ast.Ident); ok {
						// 检查是否使用了别名指向的包的函数
						if pkgPath, ok := importAliase2importPath[x.Name]; ok && pkgPath == libDaprImportPath {
							if pie.Contains(daprServerFunctions, n.Sel.Name) {
								sourceInfo.ServerSourceFile = path
							}
						}
					}
				case *ast.ImportSpec: // 记录该文件所有的导入别名和完整路径的包名
					var alias string
					if n.Name != nil {
						alias = n.Name.Name
					}
					fullPath := n.Path.Value[1 : len(n.Path.Value)-1]
					pkgName := filepath.Base(fullPath)
					if alias == "" {
						importAliase2importPath[pkgName] = fullPath
					} else {
						importAliase2importPath[alias] = fullPath
					}
				}

				return true
			})

			//	// 遍历所有的顶级声明
			//LOOP:
			//	for _, decl := range f.Decls {
			//		// 对于每个GenDecl类型的声明（包括import、const、type、var声明）
			//		if genDecl, ok := decl.(*ast.GenDecl); ok {
			//			// 仅处理类型声明
			//			if genDecl.Tok == token.TYPE {
			//				for _, spec := range genDecl.Specs {
			//					// 如果类型规范是类型别名或类型声明
			//					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			//						// 如果类型规范是结构体类型
			//						structures, ok := typeSpec.Type.(*ast.StructType)
			//						if ok {
			//							// 检查第一个field
			//							if len(structures.Fields.List) > 0 {
			//								if fmt.Sprintf("%s", structures.Fields.List[0].Type) == expr {
			//									found, _ = filepath.Rel(a.sourceDir, filepath.Dir(path))
			//									found = filepath.ToSlash(found)
			//									break LOOP
			//								}
			//							}
			//						}
			//
			//					}
			//				}
			//			}
			//		}
			//	}
		}

		return nil
	})

	return sourceInfo, nil
}

func (p moduleAstParserImpl) getProjectModuleName() (string, error) {
	// 获取根模块名
	cmdOutput, err := exec.Command("go", "list", "-m").CombinedOutput()
	if err != nil {
		return "", err
	}

	// 按换行符拆分结果
	lines := bytes.Split(cmdOutput, []byte("\n"))
	if len(lines) == 0 {
		return "", errors.New("project is not using go module or not run go list -m in project root dir")
	}

	return strings.TrimSpace(string(lines[0])), nil
}

// MonkeyPatch 修改源代码的方式匿名导入pkg, sourceFile是相对于basePath的相对路径
func (p moduleAstParserImpl) addImportModulePaths(sourceFile string, modulePaths map[ModuleKind]string) error {
	// 获取项目模块名
	projectModuleName, err := p.getProjectModuleName()
	if err != nil {
		return err
	}

	// 将源代码解析为抽象语法树（AST）
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, sourceFile, nil, 0)
	if err != nil {
		return errors.Wrapf(err, "golang ast parse file, path: %s", sourceFile)
	}

	// 记录所有已经导入的包
	allImportPaths := make(map[string]struct{})
	for _, spec := range astFile.Imports {
		allImportPaths[spec.Path.Value] = struct{}{}
	}

	// 创建新的import节点匿名插入到import声明列表
	for _, modulePath := range pie.Values(modulePaths) {
		// IMPORTANT: spec.Path.Value是带了双引号的
		checkValue := "\"" + path.Join(projectModuleName, modulePath) + "\""

		// 当patch进去的路径不存在时才加入
		if _, exists := allImportPaths[checkValue]; !exists {
			// 创建一个新的匿名ImportSpec节点
			spec := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: checkValue,
				},
				Name: ast.NewIdent("_"), // 下划线表示匿名导入
			}

			// 创建一个新的声明并插入到文件的声明列表中
			decl := &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					spec,
				},
			}

			astFile.Decls = append([]ast.Decl{decl}, astFile.Decls...)
		}
	}

	// 使用printer包将抽象语法树（AST）打印成代码
	buf := bytes.NewBufferString("")
	err = printer.Fprint(buf, fset, astFile)
	if err != nil {
		return err
	}

	// 打开文件
	file, err := os.OpenFile(sourceFile, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将新代码内容写入文件
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// 确保所有操作都已写入磁盘
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (p moduleAstParserImpl) discoverRouteItems(relInvocationModulePath string) ([]*protobuf.RouteItem, error) {
	routeItems := make([]*protobuf.RouteItem, 0)
	absModulePath := filepath.Join(p.sourceRootDir, relInvocationModulePath)
	for _, moduleInstance := range LoadInvocationModules(relInvocationModulePath) {
		routeAnnotations, err := moduleInstance.GetRouteAnnotations(absModulePath, p.handlerNameMatchers...)
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
