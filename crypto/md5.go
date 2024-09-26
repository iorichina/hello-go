package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Hex
//
// A non-constant value x can be converted to type T in any of these cases:
//
// x is assignable to T.
//
// x's type and T have identical underlying types.
//
// x's type and T are unnamed pointer types and their pointer base types have identical underlying types.
//
// x's type and T are both integer or floating point types.
//
// x's type and T are both complex types.
//
// x is an integer or a slice of bytes or runes and T is a string type.
//
// x is a string and T is a slice of bytes or runes.
//
// 规范只允许将字节 slice 或 rune 转换为字符串，而不是字节数组。
//
// 在 Go 中，数组和 slice 是不同的类型。数组的大小是类型的一部分。
func Md5Hex(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}
