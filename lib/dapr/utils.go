package dapr

import (
	"github.com/hdget/hdutils/convert"
	"github.com/hdget/hdutils/text"
)

var (
	truncateSize = 200
)

func truncate(data []byte) string {
	return text.Truncate(convert.BytesToString(data), truncateSize)
}
