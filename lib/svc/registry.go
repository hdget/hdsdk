package svc

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

var (
	registry = make(map[string]Module)
)

type testModule struct {
	Module
}

var (
	errEmptyModuleName = errors.New("empty module name")
)

func (*testModule) GetSignatureHandler(ctx context.Context, event *common.InvocationEvent) (*common.Content, error) {
	return nil, nil
}

func RegisterDaprModule(svcHolder any, app string, version int, options ...Option) error {
	moduleName := utils.GetStructName(svcHolder)
	if moduleName == "" {
		return errEmptyModuleName
	}

	err := utils.StructSetComplexField(svcHolder, (*Module)(nil), NewDaprModule(app, moduleName, version, options...))
	if err != nil {
		return errors.Wrapf(err, "set base module for: %s ", moduleName)
	}

	module, ok := svcHolder.(Module)
	if !ok {
		return errors.New("invalid module")
	}

	// 将实际的struct实例保存进去
	module.Super(svcHolder)

	// 注册handlers
	registry[module.GetName()] = module

	return nil
}

func GetRegistry() map[string]Module {
	return registry
}
