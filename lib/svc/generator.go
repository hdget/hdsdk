package svc

type Generator interface {
	Register() error          // 通过生成的源文件来注册相关信息
	Gen(srcPath string) error // 通过解析源代码来生成源文件
	Get() any                 // 获取生成的内容
}

type BaseGenerator struct {
}

func NewGenerator() Generator {
	return &BaseGenerator{}
}

func (m *BaseGenerator) Get() any {
	return nil
}

func (m *BaseGenerator) Register() error {
	return nil
}

func (m *BaseGenerator) Gen(srcPath string) error {
	//TODO implement me
	panic("implement me")
}

//
//func (m *BaseGenerator) Register() error {
//	if m.concrete == nil {
//		return errors.New("no concrete generator")
//	}
//
//	// 因为具体Generator的Register方法可能未被实现，其会导致循环调用base.Register
//	// 这里用loopCount来计数，保证无循环调用
//	if m.loopCount > 0 {
//		return fmt.Errorf("no 'Register' method implemented for concrete generator: %s", utils.Reflect().GetStructName(m.concrete))
//	}
//
//	m.loopCount += 1
//
//	return m.concrete.Register()
//	// 通过反射的方法调用Register
//	//method := reflect.ValueOf(m.concrete).MethodByName(m.GetRegisterMethodName())
//	//if !method.IsNil() {
//	//	results := method.Call(nil)
//	//	if len(results) == 0 || results[0].Type().String() != "error" {
//	//		return errors.New("invalid register results")
//	//	}
//	//	return results[0].Interface().(error)
//	//}
//	//return nil
//}
