package dapr

import "github.com/hdget/hdutils/convert"

var (
	maxTrimSize = 200
)

func trimData(data []byte) string {
	trimmed := []rune(convert.BytesToString(data))
	if len(trimmed) > maxTrimSize {
		trimmed = append(trimmed[:maxTrimSize], []rune("...")...)
	}
	return string(trimmed)
}
