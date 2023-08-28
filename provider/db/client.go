package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/elgris/sqrl"
	"github.com/hdget/hdsdk/types"
)

func (b *SqDbClient) Sq(builder squirrel.Sqlizer) types.DbClient {
	b._builder = builder
	return b
}

func (b *SqDbClient) Sqrl(builder sqrl.Sqlizer) types.DbClient {
	b._builder = builder
	return b
}
