package jsoner

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Float64 自定义类型，用于处理字符串格式的整数
type Float64 float64

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *Float64) UnmarshalJSON(data []byte) error {
	// 尝试直接解析为整数
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*i = Float64(num)
		return nil
	}

	// 尝试解析为字符串，再转为整数
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid int format: %s", data)
	}

	parsed, err := strconv.ParseFloat(str, 64)
	if err != nil {
		//return fmt.Errorf("cannot convert %q to int: %v", str, err)
	}

	*i = Float64(parsed)
	return nil
}
