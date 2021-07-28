// Package ms provides microservice ability
package ms

import (
	"github.com/hdget/sdk/types"
)

type BaseMsProvider struct {
	Items map[string]types.MicroService // 指定的微服务
}

func (p *BaseMsProvider) By(name string) types.MicroService {
	return p.Items[name]
}
