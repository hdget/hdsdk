package svc

import (
	"fmt"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
)

type Generator interface {
	GetName() string               // 获取生成器的实现结构的名字
	GetRegisterMethodName() string // 获取注册方法名
	Register() error               // 通过生成的源文件来注册相关信息
	Self(generator Generator)      // 将具体的Generator注册到base中
	Gen(srcPath string) error      // 通过解析源代码来生成源文件
}

type BaseGenerator struct {
	concrete  Generator
	loopCount int // 调用次数，用来追踪是否调的是具体Generator的函数还是base.Register函数
}

func NewGenerator() Generator {
	return &BaseGenerator{}
}

// Self 将具体实现设置为自己
func (m *BaseGenerator) Self(generator Generator) {
	m.concrete = generator
}

func (m *BaseGenerator) GetName() string {
	return utils.Reflect().GetStructName(m.concrete)
}

func (m *BaseGenerator) GetRegisterMethodName() string {
	return utils.Reflect().GetFuncName(m.Register)
}

func (m *BaseGenerator) Register() error {
	if m.concrete == nil {
		return errors.New("no concrete generator")
	}

	// 因为具体Generator的Register方法可能未被实现，其会导致循环调用base.Register
	// 这里用loopCount来计数，保证无循环调用
	if m.loopCount > 0 {
		return fmt.Errorf("no 'Register' method implemented for concrete generator: %s", utils.Reflect().GetStructName(m.concrete))
	}

	m.loopCount += 1

	return m.concrete.Register()
	// 通过反射的方法调用Register
	//method := reflect.ValueOf(m.concrete).MethodByName(m.GetRegisterMethodName())
	//if !method.IsNil() {
	//	results := method.Call(nil)
	//	if len(results) == 0 || results[0].Type().String() != "error" {
	//		return errors.New("invalid register results")
	//	}
	//	return results[0].Interface().(error)
	//}
	//return nil
}

func (m *BaseGenerator) Gen(srcPath string) error {
	//TODO implement me
	panic("implement me")
}
