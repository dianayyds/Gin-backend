package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
)

// 字符串转slice  ,分割的， 1,2,3 => [1, 2, 3]
func StringToIntSlice(str string) []int {
	ids := make([]int, 0)
	if str == "" {
		return ids
	}
	for _, pid := range strings.Split(str, ",") {
		id, _ := strconv.Atoi(pid)
		ids = append(ids, id)
	}
	return ids
}

// 字符串转slice  ,分割的， 1,2,3 => [1, 2, 3]
func StringToInt32Slice(str string) []int32 {
	ids := make([]int32, 0)
	if str == "" {
		return ids
	}
	for _, pid := range strings.Split(str, ",") {
		id, _ := strconv.Atoi(pid)
		ids = append(ids, int32(id))
	}
	return ids
}
func StringToUInt32Slice(str string) []uint32 {
	ids := make([]uint32, 0)
	if str == "" {
		return ids
	}
	for _, pid := range strings.Split(str, ",") {
		id, _ := strconv.Atoi(pid)
		ids = append(ids, uint32(id))
	}
	return ids
}

// slice转字符串  ,分割的， [1, 2, 3] => 1,2,3
func IntSliceToString(a []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func UInt32SliceToString(a []uint32) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

// slice转字符串  ,分割的， [1, 2, 3] => 1,2,3
func Int32SliceToString(a []int32) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

// 密码长度 安全程度检测，nil表示符合要求
func CheckPwd(minLength, maxLength int, pwd string) error {
	if len(pwd) > 0 {
		return nil
	}
	if len(pwd) < minLength {
		return fmt.Errorf("BAD PASSWORD: The password is shorter than %d characters", minLength)
	}
	if len(pwd) > maxLength {
		return fmt.Errorf("BAD PASSWORD: The password is logner than %d characters", maxLength)
	}

	var level int = 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}

	if level < 3 {
		return fmt.Errorf("The password does not satisfy the current policy requirements.")
	}
	return nil
}

// 当前时间的UTC时间，  time.Time  和 string
func NowUTC() time.Time {
	// 指定时区的今天此刻
	now := carbon.Now(carbon.UTC)
	return now.Carbon2Time()
}

// 某个时区的时间转换成UTC时间
func TimeZone2UTCTime(t, zoneID string) time.Time {
	if zoneID == "" {
		zoneID = carbon.UTC
	}
	return carbon.Parse(t, zoneID).Carbon2Time()
}

// UTC时间转换成某个时区的时间
func UTCTime2TimeZone(t time.Time, zoneID string) string {
	if zoneID == "" {
		zoneID = carbon.UTC
	}
	return carbon.Parse(t.Format("2006-01-02 15:04:05"), zoneID).ToDateTimeString()
}

// 判断slice里是否包含某个值
func IsInUint32Slice(ids []uint32, id uint32) bool {
	for _, k := range ids {
		if k == id {
			return true
		}
	}
	return false
}

func IsInStingSlice(target string, slice []string) bool {
	for _, fruit := range slice {
		if fruit == target {
			return true
		}
	}
	return false
}

func slicesEqual(s1, s2 []interface{}) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
