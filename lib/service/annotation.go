package service

//
//import (
//	"fmt"
//	"strings"
//)
//
//type ModuleInfo struct {
//	ModuleName  string
//	FuncName    string
//	Annotations map[string]*Annotation // annotationName->*Annotation
//	Comments    []string
//}
//
//// 注解的前缀
//const defaultAnnotationPrefix = "@hd."
//const annotationRoute = defaultAnnotationPrefix + "route"
//
//// parseComments 从函数注释中解析注解
//// 第一个值为注解map, Annotation=>Annotation value
//// 第二个值为除去注解后的注释信息
//func parseComments(comments []string, annPrefix string) (map[string]*Annotation, []string, error) {
//	plainComments := make([]string, 0)
//	annotations := make(map[string]*Annotation)
//	for _, s := range comments {
//		idxAnnotation := strings.Index(s, annPrefix)
//
//		// 找不到annotation前缀则直接添加到注释中
//		if idxAnnotation < 0 {
//			s = strings.Replace(s, "//", "", 1)
//			s = strings.TrimSpace(s)
//			plainComments = append(plainComments, s)
//			continue
//		}
//
//		// 找到匹配的注解前缀
//		// 去除掉前面的slash
//		s = s[idxAnnotation:]
//		// 尝试找到annotation name
//		fields := strings.Fields(s)
//		nameIndex := -1
//		for i, field := range fields {
//			if strings.HasPrefix(field, annPrefix) {
//				nameIndex = i
//				break
//			}
//		}
//
//		// 没找到annotation name
//		if nameIndex == -1 {
//			return nil, nil, fmt.Errorf("Annotation name not found, line: %s", s)
//		}
//
//		// 总是将找到的annotation加入到map，即保证最后一个生效
//		annotationName := fields[nameIndex]
//		annotations[annotationName] = &Annotation{
//			Name:  annotationName,
//			Value: strings.Join(fields[nameIndex+1:], ""),
//		}
//		//if _, exist := annotations[annotationName]; !exist && annotationName != "" {
//		//	annotations[annotationName] = &Annotation{
//		//		Name:  annotationName,
//		//		Value: strings.Join(fields[nameIndex+1:], ""),
//		//	}
//		//}
//
//	}
//
//	return annotations, plainComments, nil
//}

//
//// getModuleComments 获取DaprServiceHandler的注释
//func getModuleComments(moduleSrcPath string, matchFn MatchFunc) []funcComment {
//	fs := token.NewFileSet()
//	pkgAsts, err := parser.ParseDir(fs, moduleSrcPath, nil, parser.ParseComments)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	fnComments := make([]funcComment, 0)
//	for pkgName, pkgAst := range pkgAsts {
//		// 遍历每个包的每个文件
//		for filename, f := range pkgAst.Files {
//			fnInfos := make([]funcInfo, 0)
//			// 尝试获取文件中的所有函数定义，获取其函数名,receiver名,和位置信息
//			for _, decl := range f.Decls {
//				if fn, ok := decl.(*ast.FuncDecl); ok {
//					fnInfos = append(fnInfos, funcInfo{
//						FuncName:      fn.Name.Name,
//						FuncReceiver:  ast2.AstGetReceiverStructName(fn),
//						Pos:           fn.Pos(),
//						End:           fn.End(),
//						IsHandlerFunc: matchFn(fn),
//					})
//				}
//			}
//
//			// 遍历获取到的所有函数信息， 获取其Comment信息
//			for i, fn := range fnInfos {
//				if !fn.IsHandlerFunc {
//					continue
//				}
//
//				comment := funcComment{
//					Package:     pkgName,
//					File:    filename,
//					ModuleName:  fn.FuncReceiver,
//					HandlerName: fn.FuncName,
//					PlainComments:    make([]string, 0),
//				}
//
//				// 因为下面需要比较Comment的位置是否是在上一个函数之后，当前函数的开始之前
//				prevIndex := i - 1
//				if prevIndex < 0 {
//					prevIndex = 0
//				}
//
//				// 解析当前函数的注释
//				for _, cg := range f.PlainComments {
//					for _, c := range cg.List {
//						if c.Pos() >= fnInfos[prevIndex].End && c.End() <= fn.Pos {
//							comment.PlainComments = append(comment.PlainComments, c.Text)
//						}
//					}
//				}
//				fnComments = append(fnComments, comment)
//			}
//		}
//	}
//
//	return fnComments
//}
//
//// getFuncReceiverStructName 获取函数的receiver对应的结构名
//func getFuncReceiverStructName(fn *ast.FuncDecl) string {
//	if fn.Recv != nil {
//		for _, field := range fn.Recv.List {
//			if x, ok := field.Type.(*ast.StarExpr); ok {
//				return fmt.Sprintf("%v", x.X)
//			}
//			if x, ok := field.Type.(*ast.Ident); ok {
//				return x.String()
//			}
//		}
//	}
//	return ""
//}
