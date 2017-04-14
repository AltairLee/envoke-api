package common

import (
	"strconv"
	"strings"
)

func EmptyStr(s string) bool {
	return s == ""
}

func RepeatStr(s string, n int) string {
	return strings.Repeat(s, n)
}

func SplitStr(s, sep string) []string {
	return strings.Split(s, sep)
}

func FormatInt(x int64, base int) string {
	return strconv.FormatInt(x, base)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

func MustAtoi(s string) int {
	i, err := Atoi(s)
	Check(err)
	return i
}

func Itoa(x int) string {
	return strconv.Itoa(x)
}

func FormatInt64(x int64, base int) string {
	return strconv.FormatInt(x, base)
}

func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func ParseInt32(s string, base int) (int32, error) {
	x, err := strconv.ParseInt(s, base, 32)
	if err != nil {
		return 0, err
	}
	return int32(x), nil
}

func ParseInt64(s string, base int) (int64, error) {
	return strconv.ParseInt(s, base, 64)
}

func MustParseFloat64(s string) float64 {
	x, err := ParseFloat64(s)
	Check(err)
	return x
}

func MustParseInt64(s string, base int) int64 {
	x, err := ParseInt64(s, base)
	Check(err)
	return x
}

func ParseUint16(s string, base int) (int, error) {
	x, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func ParseUint32(s string, base int) (int, error) {
	x, err := strconv.ParseUint(s, base, 32)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func ParseUint64(s string, base int) (int, error) {
	x, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func MustParseUint16(s string, base int) int {
	x, err := ParseUint16(s, base)
	Check(err)
	return x
}

func MustParseUint32(s string, base int) int {
	x, err := ParseUint32(s, base)
	Check(err)
	return x
}

func MustParseUint64(s string, base int) int {
	x, err := ParseUint64(s, base)
	Check(err)
	return x
}
