package hdutils

import (
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// AstAnnotation 注释中的注解
type AstAnnotation struct {
	Name  string
	Value string
}

type AstFunction struct {
	Receiver      string
	Function      string
	Annotations   map[string]*AstAnnotation // annotationName->*AstAnnotation
	PlainComments []string                  // 去除注解后的注释
}

// astRawFunction ast解析出来的原始函数的信息
type astRawFunction struct {
	Package  string
	File     string
	Receiver string
	Name     string   // 函数名
	Comments []string // 出去注解后的其他注释内容
}

// astFuncInfo ast解析出来的函数的基本信息
type astFuncInfo struct {
	Name          string
	Receiver      string
	Pos           token.Pos
	End           token.Pos
	IsMatchedFunc bool
}

type aster interface {
	InspectFunction(srcPath string, fnParams, fnResults []string, annotationTag string) ([]*AstFunction, error)
}

type hdAst struct {
}

func AST() aster {
	return &hdAst{}
}

// InspectFunction 从源代码目录中获取fnParams和fnResults匹配的函数的信息,并解析函数对应的注解
// handlerName=>*moduleInfo
func (a *hdAst) InspectFunction(srcPath string, fnParams, fnResults []string, annotationPrefix string) ([]*AstFunction, error) {
	functions, err := a.getFunctions(srcPath, fnParams, fnResults)
	if err != nil {
		return nil, err
	}

	// 默认从moduleSrcPath目录开始解析, e,g: src/base/pkg/service
	funcInfos := make([]*AstFunction, 0)
	for _, fn := range functions {
		annotations, plainComments, err := a.parseComments(fn.Comments, annotationPrefix)
		if err != nil {
			return nil, err
		}

		funcInfos = append(funcInfos, &AstFunction{
			Receiver:      fn.Receiver,
			Function:      fn.Name,
			Annotations:   annotations,
			PlainComments: plainComments,
		})
	}

	return funcInfos, nil
}

// GetFunctionReceiverName 获取函数的receiver名字
// e,g: (*Person) hello() {}, 传入hello的ast.FuncDecl, 返回Person字符床
func (a *hdAst) getFunctionReceiverName(fn *ast.FuncDecl) string {
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

// getFunctions 获取匹配的函数信息
func (a *hdAst) getFunctions(srcPath string, fnParams, fnResults []string) ([]astRawFunction, error) {
	pkgAsts, err := parser.ParseDir(token.NewFileSet(), srcPath, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "ast parse src code comments, dir: %s", srcPath)
	}

	rawFuncs := make([]astRawFunction, 0)
	for pkgName, pkgAst := range pkgAsts {
		// 遍历每个包的每个文件
		for filename, f := range pkgAst.Files {
			fnInfos := make([]astFuncInfo, 0)
			// 尝试获取文件中的所有函数定义，获取其函数名,receiver名,和位置信息
			for _, decl := range f.Decls {
				if fnDecl, ok := decl.(*ast.FuncDecl); ok {
					funcName := fnDecl.Name.Name
					recv := a.getFunctionReceiverName(fnDecl)
					matched := a.matchFunction(fnDecl, fnParams, fnResults)
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

				rawFn := astRawFunction{
					Package:  pkgName,
					File:     filename,
					Receiver: fn.Receiver,
					Name:     fn.Name,
					Comments: make([]string, 0),
				}

				// 因为下面需要比较Comment的位置是否是在上一个函数之后，当前函数的开始之前
				prevIndex := i - 1
				if prevIndex < 0 {
					prevIndex = 0
				}

				// 解析当前函数有效的注释,尝试从最近的注释块开始判断, 如果注释组的首行的开始位置在上一个函数后，结束行的结束位置在当前函数前，
				// 则认为该注释块是该函数的有效注释块，这里需要倒序检查注释组
				// 为什么取最近的有效注释块，是因为需要忽略有空行分割的其他无效的注释块
				var validCommentGroup *ast.CommentGroup
				for i := len(f.Comments) - 1; i >= 0; i-- {
					currentCg := f.Comments[i]
					if len(currentCg.List) > 0 {
						cgFirstLine := currentCg.List[0]
						cgLastLine := currentCg.List[len(currentCg.List)-1]
						if cgFirstLine.Pos() >= fnInfos[prevIndex].End && cgLastLine.End() <= fn.Pos {
							validCommentGroup = currentCg
							break
						}
					}
				}

				if validCommentGroup != nil {
					for _, c := range validCommentGroup.List {
						rawFn.Comments = append(rawFn.Comments, c.Text)
					}
				}

				rawFuncs = append(rawFuncs, rawFn)
			}
		}
	}

	return rawFuncs, nil
}

// matchFunction 校验函数声明是否是匹配的
// e,g:
// xxxHandler(ctx context.Context, event *common.InvocationEvent) (*common.Content, error)
func (a *hdAst) matchFunction(fnDecl *ast.FuncDecl, params, results []string) bool {
	// 首先校验参数个数和返回值个数
	if len(fnDecl.Type.Params.List) != len(params) || len(fnDecl.Type.Results.List) != len(params) {
		return false
	}

	// 校验入参
	if a.countMatchedField(fnDecl.Type.Params.List, params) != len(params) {
		return false
	}

	// 校验返回结果类型
	if a.countMatchedField(fnDecl.Type.Results.List, results) != len(params) {
		return false
	}

	return true
}

func (a *hdAst) countMatchedField(fields []*ast.Field, typeNames []string) int {
	// 校验入参
	countValid := 0
	for i, field := range fields {
		var fieldName string
		switch v := field.Type.(type) {
		case *ast.Ident:
			fieldName = v.Name
		case *ast.StarExpr:
			if vv, ok := v.X.(*ast.SelectorExpr); ok {
				fieldName = "*" + a.getIndentName(vv.X) + "." + vv.Sel.Name
			} else {
				fieldName = "*" + a.getIndentName(v.X)
			}
		case *ast.SelectorExpr:
			fieldName = a.getIndentName(v.X) + "." + v.Sel.Name
		}

		// 检查参数名或者返回结果名是否与typeNames中的值相等
		if fieldName == typeNames[i] {
			countValid += 1
		}
	}
	return countValid
}

func (*hdAst) getIndentName(expr ast.Expr) string {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

// parseComments 从函数注释中解析注解
// 第一个值为注解map, AstAnnotation=>AstAnnotation value
// 第二个值为除去注解后的注释信息
func (*hdAst) parseComments(comments []string, annPrefix string) (map[string]*AstAnnotation, []string, error) {
	plainComments := make([]string, 0)
	annotations := make(map[string]*AstAnnotation)
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
			if strings.HasPrefix(strings.ToLower(field), strings.ToLower(annPrefix)) {
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
			annotations[annName] = &AstAnnotation{
				Name:  annName,
				Value: annValue,
			}
		}
		//if _, exist := annotations[annotationName]; !exist && annotationName != "" {
		//	annotations[annotationName] = &AstAnnotation{
		//		Name:  annotationName,
		//		Value: strings.Join(fields[nameIndex+1:], ""),
		//	}
		//}

	}

	return annotations, plainComments, nil
}
