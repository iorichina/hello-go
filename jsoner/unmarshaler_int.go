package jsoner

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Int 自定义类型，用于处理字符串格式的整数
type Int int

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *Int) UnmarshalJSON(data []byte) error {
	// 尝试直接解析为整数
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*i = Int(num)
		return nil
	}

	// 尝试解析为字符串，再转为整数
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid int format: %s", data)
	}

	parsed, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("cannot convert %q to int: %v", str, err)
	}

	*i = Int(parsed)
	return nil
}
