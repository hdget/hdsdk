package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var (
	DEFAULT_TIME_LOCATION = time.FixedZone("CST", 8*3600)
	ISO_DATE_TEMPLATE     = "2006-01-02"
)

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

	t, err := time.ParseInLocation(layout, value, DEFAULT_TIME_LOCATION)
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

// GetBetweenDays
// @description   获取开始日期和结束日期中间的所有日期列表
// @param beginDate string 开始时间,格式为：2020-04-01
// @param args ...string 如果指定了结束时间,则用结束时间,否则用当前时间,格式为：2020-04-01
// @return 在这段日期时间内的所有天包含起始日期 []string,如:[2020-04-01 2020-04-02 2020-04-03]
func GetBetweenDays(format, beginDate string, args ...string) ([]string, error) {
	beginTime, err := time.ParseInLocation(ISO_DATE_TEMPLATE, beginDate, DEFAULT_TIME_LOCATION)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid begin date, beginDate: %s", beginDate)
	}

	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, DEFAULT_TIME_LOCATION)
	if len(args) > 0 {
		endTime, err = time.ParseInLocation(ISO_DATE_TEMPLATE, args[0], DEFAULT_TIME_LOCATION)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid end date, endDate: %s", args[0])
		}
	}

	days := DeltaDays(beginTime, endTime)
	if days == -1 {
		// 如果结束时间小于开始时间，异常
		return nil, fmt.Errorf("invalid begin date or end date, beginDate: %s, args: %v", beginDate, args)
	}

	// 输出日期格式固定
	daySlice := make([]string, 0)
	for i := 0; i <= days; i++ {
		strNextDate := beginTime.AddDate(0, 0, i).Format(format)
		daySlice = append(daySlice, strNextDate)
	}

	return daySlice, nil
}

// DeltaDays 计算两个日期间的间隔天数
func DeltaDays(t1, t2 time.Time) int {
	if t1.Location().String() != t2.Location().String() {
		return -1
	}

	hours := t2.Sub(t1).Hours()
	if hours < 0 {
		return -1
	}
	// sub hours less than 24
	if hours < 24 {
		// may same day
		t1y, t1m, t1d := t1.Date()
		t2y, t2m, t2d := t2.Date()
		isSameDay := t1y == t2y && t1m == t2m && t1d == t2d

		if isSameDay {
			return 0
		} else {
			return 1
		}
	} else { // equal or more than 24
		if (hours/24)-float64(int(hours/24)) == 0 { // just 24's times
			return int(hours / 24)
		} else { // more than 24 hours
			return int(hours/24) + 1
		}
	}
}
