package dapr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"strings"
)

// funcInfo ast解析出来的函数的基本信息
type funcInfo struct {
	FuncName      string
	FuncReceiver  string
	Pos           token.Pos
	End           token.Pos
	IsHandlerFunc bool
}

// funcComment ast解析出来的注释
type funcComment struct {
	PkgName     string
	FileName    string
	ModuleName  string
	HandlerName string
	Comments    []string
}

// annotation 注释中的注解
type annotation struct {
	Name  string
	Value string
}

type serviceModuleHandlerAnnotation struct {
	ModuleName  string
	HandlerName string
	Annotations map[string]*annotation // annotationName->*annotation
	Comments    []string
}

// 注解的前缀
const annotationPrefix = "@hd."
const annotationRoute = annotationPrefix + "route"

// GetAnnotations 解析服务模块中的所有注解和注释
// handlerName=>*serviceModuleHandlerAnnotation
func (sm *ServiceModule) GetAnnotations(args ...string) ([]*serviceModuleHandlerAnnotation, error) {
	// 默认从src/app/pkg/service目录开始解析
	destPath := path.Join([]string{"src", sm.app, "pkg", "service"}...)
	if len(args) > 0 {
		destPath = args[0]
	}

	handlerAnnotations := make([]*serviceModuleHandlerAnnotation, 0)
	for _, fnComment := range getServiceModuleComments(destPath) {
		annotations, comments, err := parseAnnotations(fnComment.Comments)
		if err != nil {
			return nil, err
		}

		handlerAnnotations = append(handlerAnnotations, &serviceModuleHandlerAnnotation{
			ModuleName:  fnComment.ModuleName,
			HandlerName: fnComment.HandlerName,
			Annotations: annotations,
			Comments:    comments,
		})
	}

	return handlerAnnotations, nil
}

// parseAnnotations 从函数备注中解析路由备注
// annotationName => annotation
func parseAnnotations(comments []string) (map[string]*annotation, []string, error) {
	plainComments := make([]string, 0)
	annotations := make(map[string]*annotation, 0)
	for _, s := range comments {
		idxAnnotation := strings.Index(s, annotationPrefix)
		// 找不到annotation前缀则直接添加到注释中
		if idxAnnotation < 0 {
			s = strings.Replace(s, "//", "", 1)
			s = strings.TrimSpace(s)
			plainComments = append(plainComments, s)
		} else {
			// 去除掉前面的slash
			s = s[idxAnnotation:]

			// 尝试找到annotation name
			fields := strings.Fields(s)
			nameIndex := -1
			for i, field := range fields {
				if strings.HasPrefix(field, annotationPrefix) {
					nameIndex = i
					break
				}
			}

			// 没找到annotation name
			if nameIndex == -1 {
				return nil, nil, fmt.Errorf("annotation name not found, line: %s", s)
			}

			// 只将第一个找到的annotation加入到map
			annotationName := fields[nameIndex]
			if _, exist := annotations[annotationName]; !exist && annotationName != "" {

				annotations[annotationName] = &annotation{
					Name:  annotationName,
					Value: strings.Join(fields[nameIndex+1:], ""),
				}
			}
		}
	}

	return annotations, plainComments, nil
}

// getServiceModuleComments 获取DaprServiceHandler的注释
func getServiceModuleComments(destPath string) []funcComment {
	fs := token.NewFileSet()
	pkgAsts, err := parser.ParseDir(fs, destPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	fnComments := make([]funcComment, 0)
	for pkgName, pkgAst := range pkgAsts {
		// 遍历每个包的每个文件
		for filename, f := range pkgAst.Files {
			fnInfos := make([]funcInfo, 0)
			// 尝试获取文件中的所有函数定义，获取其函数名,receiver名,和位置信息
			for _, decl := range f.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					fnInfos = append(fnInfos, funcInfo{
						FuncName:      fn.Name.Name,
						FuncReceiver:  getFuncReceiverStructName(fn),
						Pos:           fn.Pos(),
						End:           fn.End(),
						IsHandlerFunc: isDaprServiceHandlerFuncType(fn.Type),
					})
				}
			}

			// 遍历获取到的所有函数信息， 获取其Comment信息
			for i, fn := range fnInfos {
				if !fn.IsHandlerFunc {
					continue
				}

				comment := funcComment{
					PkgName:     pkgName,
					FileName:    filename,
					ModuleName:  fn.FuncReceiver,
					HandlerName: fn.FuncName,
					Comments:    make([]string, 0),
				}

				// 因为下面需要比较Comment的位置是否是在上一个函数之后，当前函数的开始之前
				prevIndex := i - 1
				if prevIndex < 0 {
					prevIndex = 0
				}

				// 解析当前函数的注释
				for _, cg := range f.Comments {
					for _, c := range cg.List {
						if c.Pos() >= fnInfos[prevIndex].End && c.End() <= fn.Pos {
							comment.Comments = append(comment.Comments, c.Text)
						}
					}
				}
				fnComments = append(fnComments, comment)
			}
		}
	}

	return fnComments
}

// getFuncReceiverStructName 获取函数的receiver对应的结构名
func getFuncReceiverStructName(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		for _, field := range fn.Recv.List {
			if x, ok := field.Type.(*ast.StarExpr); ok {
				return fmt.Sprintf("%v", x.X)
			}
			if x, ok := field.Type.(*ast.Ident); ok {
				return x.String()
			}
		}
	}
	return ""
}

// isDaprServiceHandlerFuncType 校验函数声明是否是Darp的ServiceHandler,我们只关心ServiceHandler的注释
// serviceHandler的函数格式如下:
// xxxHandler(ctx context.Context, event *common.InvocationEvent) (*common.Content, error)
func isDaprServiceHandlerFuncType(fnType *ast.FuncType) bool {
	if fnType == nil || fnType.Params == nil || fnType.Results == nil {
		return false
	}

	if len(fnType.Params.List) != 2 || len(fnType.Results.List) != 2 {
		return false
	}

	// 校验入参
	countValidParams := 0
	for _, field := range fnType.Params.List {
		if x, ok := field.Type.(*ast.SelectorExpr); ok {
			if fmt.Sprintf("%s", x.X) == "context" {
				countValidParams += 1
			}
		}
		if x, ok := field.Type.(*ast.StarExpr); ok {
			if fmt.Sprintf("%s", x.X) == "&{common InvocationEvent}" {
				countValidParams += 1
			}
		}
	}
	if countValidParams != 2 {
		return false
	}

	// 校验返回值
	countValidResults := 0
	if x, ok := fnType.Results.List[0].Type.(*ast.StarExpr); ok {
		if fmt.Sprintf("%s", x.X) == "&{common Content}" {
			countValidResults += 1
		}
	}
	if x, ok := fnType.Results.List[1].Type.(*ast.Ident); ok {
		if x.String() == "error" {
			countValidResults += 1
		}
	}
	if countValidResults != 2 {
		return false
	}

	return true
}
