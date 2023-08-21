package svc

import (
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"reflect"
)

type Generator interface {
	GetName() string
	GetRegisterMethodName() string
	Register() error
	Self(generator Generator)
	Gen(srcPath string) error
}

type BaseGenerator struct {
	Concrete Generator
}

func NewGenerator() Generator {
	return &BaseGenerator{}
}

// Self 将具体实现设置为自己
func (m *BaseGenerator) Self(generator Generator) {
	m.Concrete = generator
}

func (m *BaseGenerator) GetName() string {
	return utils.Reflect().GetStructName(m.Concrete)
}

func (m *BaseGenerator) GetRegisterMethodName() string {
	return "Register"
}

func (m *BaseGenerator) Register() error {
	if m.Concrete == nil {
		return errors.New("no concrete generator")
	}

	results := reflect.ValueOf(m.Concrete).MethodByName(m.GetRegisterMethodName()).Call(nil)
	if len(results) == 0 || results[0].Type().String() != "error" {
		return errors.New("invalid register results")
	}
	return results[0].Interface().(error)
}

func (m *BaseGenerator) Gen(srcPath string) error {
	//TODO implement me
	panic("implement me")
}
