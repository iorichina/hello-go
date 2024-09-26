package number

import "strconv"

// 字符串转数字
func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, IntBase, Int64BitSize)
}
