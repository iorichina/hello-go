package number

import (
	"fmt"
	"testing"
)

func TestStrToInt64(t *testing.T) {
	v64 := "-3546343826724305832"
	if s, err := StrToInt64(v64); err == nil {
		fmt.Printf("%T, %v\n", s, s)
		if -3546343826724305832 != s {
			t.Fatalf("StrToInt64 失败 %#v\n", s)
		}
	}
	//int64, -3546343826724305832
}
