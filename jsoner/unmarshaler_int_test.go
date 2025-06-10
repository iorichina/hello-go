package jsoner

import (
	"encoding/json"
	"fmt"
	"testing"
)

// IntExample 包含兼容字段的结构体
type IntExample struct {
	ID    Int    `json:"id"`
	Name  string `json:"name"`
	Score Int    `json:"score"`
}

func TestInt_UnmarshalJSON(t *testing.T) {
	// 示例JSON数据，id为整数，score为字符串
	jsonData := `{"id": 123, "name": "Alice", "score": "98"}`

	var ex IntExample
	if err := json.Unmarshal([]byte(jsonData), &ex); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("ID: %d, Name: %s, Score: %d\n", ex.ID, ex.Name, ex.Score)
}
