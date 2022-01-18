package utils

import (
	"github.com/pkg/errors"
	"strings"
	"time"
)

var localBeijing = time.FixedZone("CST", 8*3600)

// ParseStrTime iso time string转化为时间，layout必须为 "2006-01-02 15:04:05"
func ParseStrTime(value string) (*time.Time, error) {
	tokens := strings.Split(value, " ")

	var layout string
	switch len(tokens) {
	case 0:
		layout = "2006-01-02 15:04:05"
	case 1:
		layout = "2006-01-02"
	default:
		return nil, errors.New("invalid time format, it is 'YYYY-MM-DD' or 'YYYY-MM-DD HH:MM:SS'")
	}

	t, err := time.ParseInLocation(layout, value, localBeijing)
	if err != nil {
		return nil, errors.Wrap(err, "parse in location")
	}

	return &t, nil
}

// IsValidBeginEndTime check if it is valid begin/end time
func IsValidBeginEndTime(strBeginTime, strEndTime string) error {
	// 检查date是否是有效的日期
	beginTime, err := ParseStrTime(strBeginTime)
	if err != nil {
		return errors.Wrap(err, "invalid begin time")
	}

	endTime, err := ParseStrTime(strEndTime)
	if err != nil {
		return errors.Wrap(err, "invalid end time")
	}

	if beginTime.After(*endTime) {
		return errors.New("end time should larger than begin time")
	}
	return nil
}
