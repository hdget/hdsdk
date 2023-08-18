package ast

import (
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// astFuncInfo ast解析出来的函数的基本信息
type astFuncInfo struct {
	Name          string
	Receiver      string
	Pos           token.Pos
	End           token.Pos
	IsMatchedFunc bool
}

// InspectFunctionByInOut 从源代码目录中获取fnParams和fnResults匹配的函数的信息,并解析函数对应的注解
// handlerName=>*moduleInfo
func InspectFunctionByInOut(srcPath string, fnParams, fnResults []string, annotationPrefix string) ([]*FunctionInfo, error) {
	comments, err := GetComments(srcPath, fnParams, fnResults)
	if err != nil {
		return nil, err
	}

	// 默认从moduleSrcPath目录开始解析, e,g: src/base/pkg/service
	funcInfos := make([]*FunctionInfo, 0)
	for _, astComment := range comments {
		annotations, plainComments, err := parseComments(astComment.Comments, annotationPrefix)
		if err != nil {
			return nil, err
		}

		funcInfos = append(funcInfos, &FunctionInfo{
			Receiver:      astComment.Struct,
			Function:      astComment.Function,
			Annotations:   annotations,
			PlainComments: plainComments,
		})
	}

	return funcInfos, nil
}

// GetFunctionReceiverName 获取函数的receiver名字
// e,g: (*Person) hello() {}, 传入hello的ast.FuncDecl, 返回Person字符床
func GetFunctionReceiverName(fn *ast.FuncDecl) string {
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

// GetComments 获取函数注释
func GetComments(srcPath string, fnParams, fnResults []string) ([]Comment, error) {
	pkgAsts, err := parser.ParseDir(token.NewFileSet(), srcPath, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "ast parse src code comments, dir: %s", srcPath)
	}

	comments := make([]Comment, 0)
	for pkgName, pkgAst := range pkgAsts {
		// 遍历每个包的每个文件
		for filename, f := range pkgAst.Files {
			fnInfos := make([]astFuncInfo, 0)
			// 尝试获取文件中的所有函数定义，获取其函数名,receiver名,和位置信息
			for _, decl := range f.Decls {
				if fnDecl, ok := decl.(*ast.FuncDecl); ok {
					funcName := fnDecl.Name.Name
					recv := GetFunctionReceiverName(fnDecl)
					matched := matchFunction(fnDecl, fnParams, fnResults)
					fnInfos = append(fnInfos, astFuncInfo{
						Name:          funcName,
						Receiver:      recv,
						Pos:           fnDecl.Pos(),
						End:           fnDecl.End(),
						IsMatchedFunc: matched,
					})
				}
			}

			// 遍历获取到的所有函数信息， 获取其Comment信息
			for i, fn := range fnInfos {
				if !fn.IsMatchedFunc {
					continue
				}

				comment := Comment{
					Package:  pkgName,
					File:     filename,
					Struct:   fn.Receiver,
					Function: fn.Name,
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
				comments = append(comments, comment)
			}
		}
	}

	return comments, nil
}

// matchFunction 校验函数声明是否是匹配的
// e,g:
// xxxHandler(ctx context.Context, event *common.InvocationEvent) (*common.Content, error)
func matchFunction(fnDecl *ast.FuncDecl, params, results []string) bool {
	// 首先校验参数个数和返回值个数
	if len(fnDecl.Type.Params.List) != len(params) || len(fnDecl.Type.Results.List) != len(params) {
		return false
	}

	// 校验入参
	if countMatchedField(fnDecl.Type.Params.List, params) != len(params) {
		return false
	}

	// 校验返回结果类型
	if countMatchedField(fnDecl.Type.Results.List, results) != len(params) {
		return false
	}

	return true
}

func countMatchedField(fields []*ast.Field, typeNames []string) int {
	// 校验入参
	countValid := 0
	for i, field := range fields {
		var fieldName string
		switch v := field.Type.(type) {
		case *ast.Ident:
			fieldName = v.Name
		case *ast.StarExpr:
			if vv, ok := v.X.(*ast.SelectorExpr); ok {
				fieldName = "*" + getIndentName(vv.X) + "." + vv.Sel.Name
			} else {
				fieldName = "*" + getIndentName(v.X)
			}
		case *ast.SelectorExpr:
			fieldName = getIndentName(v.X) + "." + v.Sel.Name
		}

		// 检查参数名或者返回结果名是否与typeNames中的值相等
		if fieldName == typeNames[i] {
			countValid += 1
		}
	}
	return countValid
}

func getIndentName(expr ast.Expr) string {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

// parseComments 从函数注释中解析注解
// 第一个值为注解map, Annotation=>Annotation value
// 第二个值为除去注解后的注释信息
func parseComments(comments []string, annPrefix string) (map[string]*Annotation, []string, error) {
	plainComments := make([]string, 0)
	annotations := make(map[string]*Annotation)
	for _, s := range comments {
		idxAnnotation := strings.Index(s, annPrefix)

		// 找不到annotation前缀则直接添加到注释中
		if idxAnnotation < 0 {
			s = strings.Replace(s, "//", "", 1)
			s = strings.TrimSpace(s)
			plainComments = append(plainComments, s)
			continue
		}

		// 找到匹配的注解前缀
		// 去除掉前面的slash
		s = s[idxAnnotation:]
		// 尝试找到annotation name
		fields := strings.Fields(s)
		nameIndex := -1
		for i, field := range fields {
			if strings.HasPrefix(field, annPrefix) {
				nameIndex = i
				break
			}
		}

		// 没找到annotation name
		if nameIndex == -1 {
			return nil, nil, fmt.Errorf("annotation name not found, line: %s", s)
		}

		// 总是将找到的annotation加入到map，即保证最后一个生效
		annName := fields[nameIndex]
		if annName != "" {
			// 处理注解值
			annValue := strings.Join(fields[nameIndex+1:], "")
			annValue = strings.TrimSpace(annValue)
			annotations[annName] = &Annotation{
				Name:  annName,
				Value: annValue,
			}
		}
		//if _, exist := annotations[annotationName]; !exist && annotationName != "" {
		//	annotations[annotationName] = &Annotation{
		//		Name:  annotationName,
		//		Value: strings.Join(fields[nameIndex+1:], ""),
		//	}
		//}

	}

	return annotations, plainComments, nil
}
