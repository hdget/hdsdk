package dapr

import (
	"encoding/json"
	"fmt"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"strings"
)

// ast解析出来的函数的基本信息
type funcInfo struct {
	FuncName      string
	Receiver      string
	Pos           token.Pos
	End           token.Pos
	IsHandlerFunc bool
}

// serviceHandlerComment dapr service handler的注释
type serviceHandlerComment struct {
	PkgName  string
	FileName string
	Receiver string
	Handler  string
	Comments []string
}

// RouteAnnotation 路由注解
type RouteAnnotation struct {
	App           string
	Namespace     string
	Version       int32
	Handler       string
	Endpoint      string
	Methods       []string
	CallerId      int64
	IsRawResponse bool
	IsPublic      bool
	Comments      []string
}

// 路由注解的前缀
const routeAnnotationPrefix = "@hd.route"

// ParseRouteAnnotations 尝试从函数的注释中获取路由注解
func ParseRouteAnnotations(app string, args ...string) ([]*RouteAnnotation, error) {
	// 默认从src/app/pkg/service目录开始解析
	destPath := path.Join([]string{"src", app, "pkg", "service"}...)
	if len(args) > 0 {
		destPath = args[0]
	}

	handlerRoutes := make([]*RouteAnnotation, 0)
	for _, fnComment := range getDaprServiceHandlerComments(destPath) {
		route, err := parseRouteAnnotation(app, fnComment.Receiver, fnComment.Handler, fnComment.Comments)
		if err != nil {
			return nil, err
		}
		handlerRoutes = append(handlerRoutes, route)
	}

	return handlerRoutes, nil
}

// parseRouteAnnotation 从函数备注中解析路由备注
func parseRouteAnnotation(app, receiver, handler string, comments []string) (*RouteAnnotation, error) {
	plainComments := make([]string, 0)
	annotations := make([]string, 0)
	for _, s := range comments {
		idxRoute := strings.Index(s, routeAnnotationPrefix)
		if idxRoute < 0 {
			s = strings.Replace(s, "//", "", 1)
			s = strings.TrimSpace(s)
			plainComments = append(plainComments, s)
		} else {
			routeAnnotation := strings.TrimSpace(s[idxRoute+len(routeAnnotationPrefix):])
			annotations = append(annotations, routeAnnotation)
		}
	}

	var def *RouteAnnotation
	for _, annotation := range annotations {
		// @hd.route后面为空或者@hd.route后面接空的{}都是有效的, 全部取默认值
		if annotation == "" || annotation == "{}" {
			def = &RouteAnnotation{}
		} else {
			if strings.HasPrefix(annotation, "{") && strings.HasSuffix(annotation, "}") {
				// 如果定义不为空，尝试unmarshal
				err := json.Unmarshal(utils.StringToBytes(annotation), &def)
				if err != nil {
					return nil, errors.Wrapf(err, "parse route, annotation: %s", annotation)
				}
			}
		}

		// 如果第一个注解不合法，尝试去解析下一个
		if def != nil {
			// 如果找到有效的路由定义，直接跳出循环
			version, namespace, err := parseNamespaceAndVersion(receiver)
			if err != nil {
				return nil, errors.Wrap(err, "parse namespace and version")
			}

			// 固定值
			def.App = app
			def.Version = version
			def.Namespace = namespace
			def.Handler = handler
			def.Comments = plainComments
			break
		}
	}

	return def, nil
}

// parseNamespaceAndVersion 从函数的receiver中按v<version>_<namespace>的格式解析出API版本号和命名空间
func parseNamespaceAndVersion(receiver string) (int32, string, error) {
	tokens := strings.Split(receiver, "_")
	if len(tokens) != 2 {
		return 0, "", fmt.Errorf("invalid module, it should be: v<number>_<namespace>, module: %s", receiver)
	}

	if !strings.HasPrefix(tokens[0], "v") {
		return 0, "", errors.New("invalid module version, it should be: v<number>")
	}

	version := cast.ToInt32(tokens[0][1:])
	namespace := tokens[1]
	if version == 0 || namespace == "" {
		return 0, "", fmt.Errorf("invalid namespace and version, receiver: %s", receiver)
	}

	return version, namespace, nil
}

// getDaprServiceHandlerComments 获取Dapr ServiceHandler的注释
func getDaprServiceHandlerComments(destPath string) []serviceHandlerComment {
	fs := token.NewFileSet()
	pkgAsts, err := parser.ParseDir(fs, destPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	fnComments := make([]serviceHandlerComment, 0)
	for pkgName, pkgAst := range pkgAsts {
		// 遍历每个包的每个文件
		for filename, f := range pkgAst.Files {
			fnInfos := make([]funcInfo, 0)
			// 尝试获取文件中的所有函数定义，获取其函数名,receiver名,和位置信息
			for _, decl := range f.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					fnInfos = append(fnInfos, funcInfo{
						FuncName:      fn.Name.Name,
						Receiver:      getFuncReceiverStructName(fn),
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

				comment := serviceHandlerComment{
					PkgName:  pkgName,
					FileName: filename,
					Receiver: fn.Receiver,
					Handler:  fn.FuncName,
					Comments: make([]string, 0),
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
