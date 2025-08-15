package utime

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	TIME_LAYOUT_WINDOWS = "2006-01-02 15-04-05"
	TIME_LAYOUT         = "2006-01-02 15:04:05"
	TIME_LAYOUT_DATE    = "2006-01-02"
	TIME_LAYOUT_Day     = "20060102"
)

var WeekDayMap = map[string]int{
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
	"Sunday":    0,
}

var MonthMap = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

const (
	SECOND     int64 = 1
	MINUTE_SEC       = 60 * SECOND
	HOUR_SEC         = 60 * MINUTE_SEC
	DAY_SEC          = 24 * HOUR_SEC
)

var OffsetTime int64 = 0

// 设置补偿时间
func SetOffsetTime(t int64) {
	OffsetTime = t
}

// 标准时间转时间戳
func Date2Unix(Y int, M int, D int, H int, I int, S int) int64 {
	str := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", Y, M, D, H, I, S)
	return DateStr2Unix(str)
}

// 错误返回-1
func DateStr2Unix(str string) int64 {
	return DataStr2UnixByLayOut(str, "2006-01-02 15:04:05")
}

// 错误返回-1
func DataStr2UnixByLayOut(dateStr, layout string) int64 {
	_, offset := GetNow().Zone()
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return -1
	}
	return t.Unix() - int64(offset)
}

func DataStr2Time(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return t
	}
	return t
}

func DataStr2LocalTime(dateStr string) time.Time {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc)
	if err != nil {
		return t
	}
	return t
}

func TimeToStr(time time.Time) string {
	return time.Format(TIME_LAYOUT)
}

func TimeToStrWindows(time time.Time) string {
	return time.Format(TIME_LAYOUT_WINDOWS)
}

// 获取当前秒时间戳
func GetNowUnix() int64 {
	return GetNow().Unix()
}

// 获取纳秒时间戳
func GetNowUnixNano() int64 {
	return GetNow().UnixNano()
}

// 获取毫秒时间戳
func GetNowUnixM() int64 {
	return GetNow().UnixNano() / 1e6
}

// 获取当前时间
func GetNow() time.Time {
	now := time.Now().Add(time.Duration(OffsetTime) * time.Second)
	return now
}

// 获取当天零点
func GetDayDot() int64 {
	timeStr := GetNow().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	_, offset := GetNow().Zone()
	timeNumber := t.Unix() - int64(offset)
	return timeNumber
}

// 获取当天剩余秒数
func GetTodayRemainSecond() int64 {
	now := GetNow()
	todayDot := GetDayDot()
	return DAY_SEC - (now.Unix() - todayDot)
}

// 获取本周第一天 (星期一) format:2006-01-02 15:04:05
func GetWeekFirstTime() time.Time {
	now := GetNow()
	year, month, day := now.Date()
	todayBegin := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	weekday := int(todayBegin.Weekday()) // 0 ~ 6

	if weekday < 1 {
		weekday = weekday + 6
	} else {
		weekday = weekday - 1
	}
	return todayBegin.AddDate(0, 0, -weekday)
}

// 获取本月第一天 format:2006-01-02 15:04:05
func GetMonthFirstTime() time.Time {
	now := GetNow()
	year, month, _ := now.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
}

// 获取时间对应的本月第一天 format:2006-01-02 15:04:05
func GetMonthFirstTimeByTime(getTime time.Time) time.Time {
	year, month, _ := getTime.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, getTime.Location())
}

// 获取月开始
func GetMonthStartTime(i int) time.Time {
	now := GetNow()
	year, month, _ := now.Date()
	return time.Date(year, month+time.Month(i), 1, 0, 0, 0, 0, now.Location())
}

// 获取上月第一天 format:2006-01-02 15:04:05
func GetLastMonthFirstTime() time.Time {
	return GetMonthFirstTime().AddDate(0, -1, 0)
}

func WeekDay() int {
	wd := GetNow().Weekday().String()
	return WeekDayMap[wd]
}

func WeekOfYear() (Year int, Week int) {
	w := WeekDay()
	m := 0
	if w == 0 {
		m = 6
	} else {
		m = w - 1
	}
	d1 := GetNowUnix() - int64(m*86400)
	tm := time.Unix(d1, 0)
	return tm.ISOWeek()
}

// 获取前s周周一
func WeekOfMondayBefore(s int) int32 {
	w := WeekDay()
	m := 0
	if w == 0 {
		m = 6
	} else {
		m = w - 1
	}
	m = m + s*7
	d1 := GetNowUnix() - int64(m*86400)
	tm := time.Unix(d1, 0)
	y, M, d := tm.Date()
	m = MonthMap[M.String()]
	return int32(y*10000 + m*100 + d)
}
func WeekOfMonday() int32 {
	w := WeekDay()
	m := 0
	if w == 0 {
		m = 6
	} else {
		m = w - 1
	}
	d1 := GetNowUnix() - int64(m*86400)
	tm := time.Unix(d1, 0)
	y, M, d := tm.Date()
	m = MonthMap[M.String()]
	return int32(y*10000 + m*100 + d)
}

func IsSameMonth(time1 int64, time2 int64) bool {
	tm1 := time.Unix(time1, 0)
	tm2 := time.Unix(time2, 0)
	y1, m1, _ := tm1.Date()
	y2, m2, _ := tm2.Date()
	return y1 == y2 && m1 == m2
}

// 是否同一天
func IsSameDay(time1 int64, time2 int64) bool {
	tm1 := time.Unix(time1, 0)
	tm2 := time.Unix(time2, 0)
	y1, m1, d1 := tm1.Date()
	y2, m2, d2 := tm2.Date()
	fmt.Println(y1, m1, d1, y2, m2, d2)
	return y1 == y2 && m1 == m2 && d1 == d2
}

// 下一个时间点，小时
func NextHour(h int64, m int64, s int64) int {
	n2 := int(h*3600 + m*60 + s)
	return NextHourSec(n2)
}
func NextHourSec(n2 int) int {
	th := GetNow().Hour()
	tm := GetNow().Minute()
	ts := GetNow().Second()
	n1 := th*3600 + tm*60 + ts
	if n1 >= n2 {
		return n2 + 86400 - n1
	} else {
		return n2 - n1
	}
}

// 下一个时间点，分钟
func NextMin(m int64, s int64) int {
	n2 := int(m*60 + s)
	return NextMinSec(n2)
}

// 下一个时间点，分钟
func NextMinSec(n2 int) int {
	tm := GetNow().Minute()
	ts := GetNow().Second()
	n1 := tm*60 + ts
	if n1 >= n2 {
		return n2 + 3600 - n1
	} else {
		return n2 - n1
	}
}

// 间隔天数
func IntervalDay(t1, t2 int64) int64 {
	if t1 > t2 {
		t3 := t1
		t1 = t2
		t2 = t3
	}
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	y1 := tm1.Year()
	m1 := tm1.Month().String()
	d1 := tm1.Day()
	y2 := tm2.Year()
	m2 := tm2.Month().String()
	d2 := tm2.Day()
	u1 := Date2Unix(y1, MonthMap[m1], d1, 0, 0, 0)
	u2 := Date2Unix(y2, MonthMap[m2], d2, 0, 0, 0)
	f := math.Floor((float64(u2) - float64(u1)) / 86400)
	return int64(f)
}

func GetDateTime() (y int, m int, d int, h int, i int, s int) {
	now := GetNow()
	y, M, d := now.Date()
	m = MonthMap[M.String()]
	h = now.Hour()
	i = now.Minute()
	s = now.Second()
	return
}

func ToDateTime(unixTime int64) (y int, m int, d int, h int, i int, s int) {
	tm1 := time.Unix(unixTime, 0)
	y, M, d := tm1.Date()
	m = MonthMap[M.String()]
	h = tm1.Hour()
	i = tm1.Minute()
	s = tm1.Second()
	return
}

// 返回两个时间的区间
func TimeSection(startTime, endTime time.Time) string {
	str := "0天"
	second := endTime.Sub(startTime).Milliseconds() / 1000
	if second >= 86400 {
		day := second / 86400
		surplusSecond := second - 86400*day //多余秒数
		str = fmt.Sprintf("%v天", day)
		str = calculateHours(surplusSecond, str)
	} else {
		str = calculateHours(second, str)
	}
	return str
}

func calculateHours(surplusSecond int64, str string) string {
	if surplusSecond >= 3600 {
		h := surplusSecond / 3600
		str += fmt.Sprintf("%v小时", h)
		surplusSecond = surplusSecond - 3600*h
		str = calculateMinute(surplusSecond, str)
	} else {
		str += "0小时"
		str = calculateMinute(surplusSecond, str)
	}
	return str
}

func calculateMinute(surplusSecond int64, str string) string {
	if surplusSecond >= 60 {
		m := surplusSecond / 60
		str += fmt.Sprintf("%v分钟%v秒", m, surplusSecond-60*m)
	} else {
		str += fmt.Sprintf("0分钟%v秒", surplusSecond)
	}
	return str
}

// 返回[yyyyMM, yyyyMM]由近到远
func GetYYYYMMList(begin int64, end int64) []string {
	if begin > end {
		begin, end = end, begin
	}
	by, bm, _ := time.Unix(begin, 0).Date()
	ey, em, _ := time.Unix(end, 0).Date()
	res := make([]string, 0)
	for {
		if ey == by && em == bm {
			now := fmt.Sprintf("%d%02d", ey, em)
			res = append(res, now)
			break
		}
		now := fmt.Sprintf("%d%02d", ey, em)

		res = append(res, now)
		em--
		if em == 0 {
			em = 12
			ey--
		}
	}
	return res
}

// 标准时间转时间戳
func Unix2YearMonth(timestamp int64) int32 {
	tm := time.Unix(timestamp, 0)
	y, M, _ := tm.Date()
	m := MonthMap[M.String()]
	return int32(y)*100 + int32(m)
}

func Unix2YearMonthList(st, et int64) []int32 {
	if st > et {
		et, st = st, et
	}
	endMonth := Unix2YearMonth(et)
	startMonth := Unix2YearMonth(st)
	list := make([]int32, 0)
	list = append(list, startMonth)
	if endMonth == startMonth {
		return list
	}
	for i := 1; i <= 100; i++ {
		cur := startMonth % 100
		if cur == 12 {
			startMonth = startMonth + 100 - 11
			list = append(list, startMonth)
		} else {
			startMonth = startMonth + 1
			list = append(list, startMonth)
		}
		if startMonth == endMonth {
			return list
		}
	}
	return list
}

// st < et
func DateBefore(st, et time.Time) bool {
	sy, sm, sd := st.Date()
	ey, em, ed := et.In(st.Location()).Date()
	ns := time.Date(sy, sm, sd, 0, 0, 0, 0, st.Location())
	ne := time.Date(ey, em, ed, 0, 0, 0, 0, st.Location())
	return ns.Before(ne)
}

func DateZero(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func ToDay() time.Time {
	now := GetNow()
	return DateZero(now)
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的23:59:59点时间
func GetDayEndTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取周一、周日的时间
func GetFirstDateOfWeek(d time.Time) (time.Time, time.Time) {
	offsetFirst := int(time.Monday - d.Weekday())
	if offsetFirst > 0 {
		offsetFirst = -6
	}
	offsetLast := 7 - int(d.Weekday())
	if offsetLast == 7 {
		offsetLast = 0
	}
	weekStartDate := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offsetFirst)
	weekEndDate := time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, time.Local).AddDate(0, 0, offsetLast)
	return weekStartDate, weekEndDate
}

// 根据年份和月份获取当月的天数
func getDays(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}

// 获取两个时间区间的年份加月份加天数 用于分表查询 user_20200101 不支持跨年度
// 时间格式为startTime 2020-01-30 00:00:00 endTime 2020-05-30 00:00:00
func GetSectionDay(startTime, endTime string) []string {
	startDate, _ := time.ParseInLocation(TIME_LAYOUT, startTime, time.Local)
	endDate, _ := time.ParseInLocation(TIME_LAYOUT, endTime, time.Local)
	//获取年
	startYear := startDate.Year()
	daySlice := make([]string, 0)
	//获取月份对应的数字
	startMonth := MonthMap[startDate.Month().String()]
	endMonth := MonthMap[endDate.Month().String()]

	startDay := startDate.Day()
	endDay := endDate.Day()
	//起始月份到结束月份 在拼接年份即可
	if startMonth == endMonth {
		for k := startMonth; k <= endMonth; k++ {
			for j := startDay; j <= endDay; j++ {
				temp := fmt.Sprintf("%d", startYear)
				//月份小于10
				if k < 10 { //小于10  格式为202001 而不是20201
					temp += fmt.Sprintf("%d%d", 0, k)
				} else {
					temp += fmt.Sprintf("%d", k)
				}
				//日份小于10
				if j < 10 { //小于10  格式为202001 而不是2020101 而不是202011
					temp += fmt.Sprintf("%d%d", 0, j)
				} else {
					temp += fmt.Sprintf("%d", j)
				}
				daySlice = append(daySlice, temp)
			}
		}
	} else {
		for k := startMonth; k <= endMonth; k++ {
			days := getDays(startYear, startMonth)
			y := 1
			if k == endMonth {
				days = endDay
			} else {
				y = startDay
			}
			for j := y; j <= days; j++ {
				temp := fmt.Sprintf("%d", startYear)
				//月份小于10
				if k < 10 { //小于10  格式为202001 而不是20201
					temp += fmt.Sprintf("%d%d", 0, k)
				} else {
					temp += fmt.Sprintf("%d", k)
				}
				//日份小于10
				if j < 10 { //小于10  格式为202001 而不是2020101 而不是202011
					temp += fmt.Sprintf("%d%d", 0, j)
				} else {
					temp += fmt.Sprintf("%d", j)
				}
				daySlice = append(daySlice, temp)
			}
		}
	}
	return daySlice
}

// 字符串转time 格式为2006-01-02 15:04:05
func StrToTime(str string) time.Time {
	endDate, _ := time.ParseInLocation(TIME_LAYOUT, str, time.Local)
	return endDate
}

// 获取两个时间区间的年份加月份 用于分表查询 user_202001 可跨年度
// 时间格式为startTime 2020-01-30 00:00:00 endTime 2021-05-30 00:00:00
func GetSectionMonth(startTime, endTime string) []string {
	startDate, _ := time.ParseInLocation(TIME_LAYOUT, startTime, time.Local)
	endDate, _ := time.ParseInLocation(TIME_LAYOUT, endTime, time.Local)
	//获取年
	startYear := startDate.Year()
	endYear := endDate.Year()
	monthSlice := make([]string, 0)
	//获取月份对应的数字
	monthStart := MonthMap[startDate.Month().String()]
	monthEnd := MonthMap[endDate.Month().String()]
	//当没有跨年度时,简单的for搞定
	if startYear == endYear {
		//起始月份到结束月份 在拼接年份即可
		for k := monthStart; k <= monthEnd; k++ {
			temp := fmt.Sprintf("%d%d", startYear, k)
			if k < 10 { //小于10  格式为202001 而不是20201
				temp = fmt.Sprintf("%d%d%d", startYear, 0, k)
			}
			monthSlice = append(monthSlice, temp)
		}
	} else if startYear < endYear {
		//循环年度
		for j := startYear; j <= endYear; j++ {
			year := 0
			if j == startYear {
				year = startYear
			} else {
				year = j
			}
			//当是最后一年时、获取最后一年的月份作为for结束条件
			if j == endYear {
				for k := 1; k <= monthEnd; k++ {
					temp1 := fmt.Sprintf("%d%d", year, k)
					if k < 10 {
						temp1 = fmt.Sprintf("%d%d%d", year, 0, k)
					}
					monthSlice = append(monthSlice, temp1)
				}
			} else {
				if j == startYear {
					for k := monthStart; k <= 12; k++ {
						temp1 := fmt.Sprintf("%d%d", year, k)
						if k < 10 {
							temp1 = fmt.Sprintf("%d%d%d", year, 0, k)
						}
						monthSlice = append(monthSlice, temp1)
					}
				} else {
					for k := 1; k <= 12; k++ {
						temp1 := fmt.Sprintf("%d%d", year, k)
						if k < 10 {
							temp1 = fmt.Sprintf("%d%d%d", year, 0, k)
						}
						monthSlice = append(monthSlice, temp1)
					}
				}
			}
		}
	}
	return monthSlice
}

// 获取当前时间最近三个月的 分表数组  格式为 202001  202002 202003
func GetLatelyMarch() []string {
	m := make([]string, 0)
	t := time.Now()
	i := MonthMap[t.Month().String()]
	year := t.Year()
	t1year := t.Year()
	t1i := i
	if i == 2 {
		i = 14
		year--
	}
	if i == 1 {
		i = 13
		year--
		t1year--
		t1i = 13
	}
	t2 := i - 2
	t1 := t1i - 1
	if t2 < 10 {
		s := strconv.Itoa(year) + "0" + strconv.Itoa(t2)
		m = append(m, s)
	} else {
		s := strconv.Itoa(year) + strconv.Itoa(t2)
		m = append(m, s)
	}
	if t1 < 10 {
		s := strconv.Itoa(t1year) + "0" + strconv.Itoa(t1)
		m = append(m, s)
	} else {
		s := strconv.Itoa(t1year) + strconv.Itoa(t1)
		m = append(m, s)
	}
	m = append(m, t.Format("200601"))
	return m
}

// 随机睡眠 随机时间 单位秒
func RandomSleep(least int, rand_range int) {
	rand.New(rand.NewSource(GetNowUnixNano()))
	// 生成least到least+range之间的随机数
	n := rand.Intn(rand_range) + least
	time.Sleep(time.Duration(n) * time.Second)
}
