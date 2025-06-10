package jsoner

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Int64 自定义类型，用于处理字符串格式的整数
type Int64 int64

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *Int64) UnmarshalJSON(data []byte) error {
	// 尝试直接解析为整数
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*i = Int64(num)
		return nil
	}

	// 尝试解析为字符串，再转为整数
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid int64 format: %s", data)
	}

	parsed, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot convert %q to int64: %v", str, err)
	}

	*i = Int64(parsed)
	return nil
}
