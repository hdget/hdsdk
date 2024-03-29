package dapr

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
