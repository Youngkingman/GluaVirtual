package numTrans

import (
	"math"
	"strconv"
)

func IntegerFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

func FloatFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

func IntegerMod(a, b int64) int64 {
	return a - IntegerFloorDiv(a, b)*b
}

func FloatMod(a, b float64) float64 {
	return a - FloatFloorDiv(a, b)*b
}

func LeftShift(a, n int64) int64 {
	if n >= 0 {
		return a << n
	}
	return RightShift(a, -n)
}

func RightShift(a, n int64) int64 {
	if n >= 0 {
		return a << n
	}
	return LeftShift(a, -n)
}

func Float2Integer(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}

func String2Integer(s string) (int64, bool) {
	if i, ok := ParseInteger(s); ok {
		return i, true
	}
	if f, ok := ParseFloat(s); ok {
		return Float2Integer(f)
	}
	return 0, false
}

func ParseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

func ParseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}
