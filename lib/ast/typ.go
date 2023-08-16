package ast

import (
	"go/ast"
)

type MatchFunc func(fnDecl *ast.FuncDecl) bool // 校验是否是需要匹配的函数

// Annotation 注释中的注解
type Annotation struct {
	Name  string
	Value string
}

// Comment ast解析出来的注释
type Comment struct {
	Package  string
	File     string
	Struct   string
	Function string
	Comments []string
}

type FunctionInfo struct {
	Receiver      string
	Function      string
	Annotations   map[string]*Annotation // annotationName->*Annotation
	PlainComments []string               // 去除注解后的注释
}
