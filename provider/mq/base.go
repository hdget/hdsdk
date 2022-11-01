package mq

import (
	"hdsdk/types"
)

type BaseMqProvider struct {
	Default types.Mq            //
	Items   map[string]types.Mq // 额外数据库
}

func (p *BaseMqProvider) My() types.Mq {
	return p.Default
}

func (p *BaseMqProvider) By(name string) types.Mq {
	return p.Items[name]
}
