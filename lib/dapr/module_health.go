package dapr

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	reflectUtils "github.com/hdget/hdutils/reflect"
	"github.com/pkg/errors"
)

type HealthModule interface {
	moduler
	GetHandler() common.HealthCheckHandler
}

type healthModuleImpl struct {
	moduler
	fn HealthCheckFunction
}

type HealthCheckFunction func(context.Context) error

var (
	_                        HealthModule = (*healthModuleImpl)(nil)
	EmptyHealthCheckFunction              = func(ctx context.Context) (err error) { return nil }
)

// NewHealthModule 健康模块
func NewHealthModule(app string, moduleObject HealthModule, fn HealthCheckFunction) error {
	// 首先实例化module
	module, err := AsHealthModule(app, moduleObject, fn)
	if err != nil {
		return err
	}

	// 最后注册module
	registerHealthModule(module)

	return nil
}

// AsHealthModule 将一个any类型的结构体转换成HealthModule
//
// e,g:
//
//		type v1_test struct {
//		  HealthModule
//		}
//
//		 v := &v1_test{}
//		 im, err := AsHealthModule("app",v)
//	     if err != nil {
//	      ...
//	     }
func AsHealthModule(app string, moduleObject any, fn HealthCheckFunction) (HealthModule, error) {
	m, err := newModule(app, moduleObject)
	if err != nil {
		return nil, err
	}

	moduleInstance := &healthModuleImpl{
		moduler: m,
		fn:      fn,
	}

	// 初始化module
	err = reflectUtils.StructSet(moduleObject, (*HealthModule)(nil), moduleInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "install module: %+v", m)
	}

	module, ok := moduleObject.(HealthModule)
	if !ok {
		return nil, errors.New("invalid health module")
	}

	return module, nil
}

func (m *healthModuleImpl) GetHandler() common.HealthCheckHandler {
	if m.fn == nil {
		return EmptyHealthCheckFunction
	}

	return func(ctx context.Context) error {
		return m.fn(ctx)
	}
}
