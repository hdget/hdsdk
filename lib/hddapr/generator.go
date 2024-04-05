package hddapr

type Generator interface {
	Gen() error // 通过解析源代码来生成源文件
	Get() any   // 获取生成的内容
}

type baseGeneratorImpl struct {
}

func (b baseGeneratorImpl) Register() error {
	//TODO implement me
	panic("implement me")
}

func (b baseGeneratorImpl) Gen() error {
	//TODO implement me
	panic("implement me")
}

func (b baseGeneratorImpl) Get() any {
	//TODO implement me
	panic("implement me")
}
