package hdutils

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var (
	DefaultTimeLocation = time.FixedZone("CST", 8*3600)
	LayoutIsoDate       = "2006-01-02"
)

// ParseStrTime iso time string转化为时间，layout必须为 "2006-01-02 15:04:05"
func ParseStrTime(value string) (*time.Time, error) {
	tokens := strings.Split(value, " ")

	var layout string
	switch len(tokens) {
	case 1:
		layout = "2006-01-02"
	case 2:
		layout = "2006-01-02 15:04:05"
	default:
		return nil, errors.New("invalid time format, it is 'YYYY-MM-DD' or 'YYYY-MM-DD HH:MM:SS'")
	}

	t, err := time.ParseInLocation(layout, value, DefaultTimeLocation)
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
	beginTime, err := time.ParseInLocation(LayoutIsoDate, beginDate, DefaultTimeLocation)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid begin date, beginDate: %s", beginDate)
	}

	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, DefaultTimeLocation)
	if len(args) > 0 {
		endTime, err = time.ParseInLocation(LayoutIsoDate, args[0], DefaultTimeLocation)
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

// GetBeginUnixTS 给出字符串的日期，例如2006-01或者2006-01-02, 返回对应的时间戳
func GetBeginUnixTS(beginDate string) int64 {
	var t time.Time
	tokens := strings.Split(beginDate, "-")
	switch len(tokens) {
	case 1:
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-01-01 00:00:00", tokens[0]), DefaultTimeLocation)
	case 2:
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-01 00:00:00", tokens[0], tokens[1]), DefaultTimeLocation)
	case 3:
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s 00:00:00", tokens[0], tokens[1], tokens[2]), DefaultTimeLocation)
	}
	return t.Unix()
}

func GetEndUnixTS(endDate string) int64 {
	var t time.Time
	tokens := strings.Split(endDate, "-")
	switch len(tokens) {
	// 提供了年
	case 1:
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-12-31 23:59:59", tokens[0]), DefaultTimeLocation)
	// 提供了年和月
	case 2:
		// 获取该月的第一天
		firstDay, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-01 00:00:00", tokens[0], tokens[1]), DefaultTimeLocation)
		// 获取该月的最后一天
		lastDay := firstDay.AddDate(0, 1, -1)
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%d 23:59:59", tokens[0], tokens[1], lastDay.Day()), DefaultTimeLocation)
	// 提供了年、月和日
	case 3:
		t, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s 23:59:59", tokens[0], tokens[1], tokens[2]), DefaultTimeLocation)
	}
	return t.Unix()
}

// GetMonthBeginTime 获取当前时间所在月份指定月份的第一天的开始时间，
// 即某月第一天的0点
// 如果nmonth=-1, 则是上一个月的第一天的00:00:00,
// 如果nmonth=0, 则是本月的第一天的00:00:00
// 如果nmonth=1, 则是下个月的第一天的00:00:00
func GetMonthBeginTime(nmonth int) time.Time {
	now := time.Now()
	theFirstDay := now.AddDate(0, 0+nmonth, -now.Day()+1)
	theFirstDayTime := time.Date(theFirstDay.Year(), theFirstDay.Month(), theFirstDay.Day(), 0, 0, 0, 0, DefaultTimeLocation)
	return theFirstDayTime
}

// GetMonthEndTime 获取当前时间的指定月份的最后一天的23:59:59
// 即某月最后一天的23:59:59
// 如果nmonth=-1, 则是上一个月的最后一天的23:59:59
// 如果nmonth=0, 则是本月的最后一天的23:59:59
// 如果nmonth=1, 则是下个月的最后一天的23:59:59
func GetMonthEndTime(nmonth int) time.Time {
	// 获取指定月的下一个月的第一天00:00:00
	nextMonthFirstDay := GetMonthBeginTime(nmonth + 1)
	// 下一个月的第一天倒退一天就是上个月的最后一天
	theLastDay := nextMonthFirstDay.AddDate(0, 0, -1)
	theLastDayTime := time.Date(theLastDay.Year(), theLastDay.Month(), theLastDay.Day(), 23, 59, 59, 0, DefaultTimeLocation)
	return theLastDayTime
}

// GetYearBeginTime 获取当前时间所在年份指定年的第一天的开始时间，
// 即某年第一天的0点
// 如果nyear=-1, 则是上一年的第一天的00:00:00,
// 如果nyear=0, 则是本年的第一天的00:00:00
// 如果nyear=1, 则是下一年的第一天的00:00:00
func GetYearBeginTime(nyear int) time.Time {
	now := time.Now()
	theFirstDay := now.AddDate(0+nyear, -int(now.Month())+1, -now.Day()+1)
	theFirstDayTime := time.Date(theFirstDay.Year(), theFirstDay.Month(), theFirstDay.Day(), 0, 0, 0, 0, DefaultTimeLocation)
	return theFirstDayTime
}

// GetYearEndTime 获取当前时间的指定年份的最后一天的23:59:59
// 即某年最后一天的23:59:59
// 如果nyear=-1, 则是上一个年的最后一天的23:59:59
// 如果nyear=0, 则是本年的最后一天的23:59:59
// 如果nyear=1, 则是下一年的最后一天的23:59:59
func GetYearEndTime(nyear int) time.Time {
	// 获取指定年的下一年的第一天00:00:00
	nextMonthFirstDay := GetYearBeginTime(nyear + 1)
	// 下一年的第一年倒退一天就是上一年的最后一天
	theLastDay := nextMonthFirstDay.AddDate(0, 0, -1)
	theLastDayTime := time.Date(theLastDay.Year(), theLastDay.Month(), theLastDay.Day(), 23, 59, 59, 0, DefaultTimeLocation)
	return theLastDayTime
}

// GetDayEndTime 获取当前时间n天后最后一秒的时间, 当前时间后n天后的日期23:59:59时间戳
// ndays: -1表示前一天，0表示今天，1表示后一天
func GetDayEndTime(ndays int) time.Time {
	now := time.Now()

	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, DefaultTimeLocation).AddDate(0, 0, ndays)
	return endTime
}

// GetDayEndTimeSince 获取从某个时间n天后最后一秒的时间, 指定时间之前的日期23:59:59时间戳
// ndays: -1表示前一天，0表示今天，1表示后一天
func GetDayEndTimeSince(ts int64, ndays int) time.Time {
	tm := time.Unix(ts, 0)
	endTime := time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, DefaultTimeLocation).AddDate(0, 0, ndays)
	return endTime
}

// GetDayBeginTimeSince 获取从某个时间n天后第一秒的时间, 指定时间之前的日期00:00:00时间戳
// ndays: -1表示前一天，0表示今天，1表示后一天
func GetDayBeginTimeSince(ts int64, ndays int) time.Time {
	tm := time.Unix(ts, 0)
	beginTime := time.Date(tm.Year(), tm.Month(), tm.Day(), 00, 00, 00, 0, DefaultTimeLocation).AddDate(0, 0, ndays)
	return beginTime
}

// GetMonthBeginTimeSince 获取从某个时间n个月后第一天第一秒的时间
// nmonth: -1表示前一个月，0表示本月，1表示后一个月
func GetMonthBeginTimeSince(ts int64, nmonth int) time.Time {
	tm := time.Unix(ts, 0)
	theFirstDay := tm.AddDate(0, 0+nmonth, -tm.Day()+1)
	theFirstDayTime := time.Date(theFirstDay.Year(), theFirstDay.Month(), theFirstDay.Day(), 0, 0, 0, 0, DefaultTimeLocation)
	return theFirstDayTime
}

// GetMonthEndTimeSince 获取从某个时间n个月最后一天最后一秒的时间, 指定时间之前的日期23:59:59时间戳
// nmonth: -1表示前一个月，0表示本月，1表示后一个月
func GetMonthEndTimeSince(ts int64, nmonth int) time.Time {
	// 获取指定月的下一个月的第一天00:00:00
	nextMonthFirstDay := GetMonthBeginTimeSince(ts, nmonth+1)
	// 下一个月的第一天倒退一天就是上个月的最后一天
	theLastDay := nextMonthFirstDay.AddDate(0, 0, -1)
	theLastDayTime := time.Date(theLastDay.Year(), theLastDay.Month(), theLastDay.Day(), 23, 59, 59, 0, DefaultTimeLocation)
	return theLastDayTime
}

// Get1stDayOfWeek 获取本周第一天
func Get1stDayOfWeek() string {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate.Format("2006-01-02")
}

func FromUnixTime(ts int64, format string) string {
	if ts <= 0 {
		return ""
	}
	tm := time.Unix(ts, 0)
	return tm.Format(format)
}
