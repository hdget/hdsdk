package service

import (
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	errEmptyModuleName = errors.New("empty module name")
)

func InitializeModule(app string, m any, args ...any) (string, error) {
	moduleName := utils.GetStructName(m)
	if moduleName == "" {
		return "", errEmptyModuleName
	}

	version := 0
	if len(args) > 0 {
		version = cast.ToInt(args[0])
	}

	err := utils.StructSet(m, &BaseModule{}, NewBaseModule(app, moduleName, version))
	if err != nil {
		return "", errors.Wrapf(err, "init base module for: %s ", moduleName)
	}

	return moduleName, nil
}

func Recover(app string) {
	if r := recover(); r != nil {
		utils.RecordErrorStack(app)
	}
}
