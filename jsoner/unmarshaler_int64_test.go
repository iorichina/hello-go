package jsoner

import (
	"encoding/json"
	"fmt"
	"testing"
)

// IntExample 包含兼容字段的结构体
type Int64Example struct {
	ID    Int64  `json:"id"`
	Name  string `json:"name"`
	Score Int64  `json:"score"`
}

func TestInt64_UnmarshalJSON(t *testing.T) {
	// 示例JSON数据，id为整数，score为字符串
	jsonData := `{"id": 123, "name": "Alice", "score": "98"}`

	var ex Int64Example
	if err := json.Unmarshal([]byte(jsonData), &ex); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("ID: %d, Name: %s, Score: %d\n", ex.ID, ex.Name, ex.Score)
}
