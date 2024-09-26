package number

import (
	"fmt"
	"hello-go/logger"
	"math/big"
)

// string 安全转 float64
func StrToFloat64(s string) (float64, bool) {
	//panic("unreachable")
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("StrToFloat64 error %#v", err)
		}
	}()
	v, b := new(big.Float).SetString(s)
	if !b {
		return 0, false
	}
	f64, _ := v.Float64()
	return f64, true
}

// float32 安全转 float64
func Float32To64(f float32) (float64, bool) {
	//panic("unreachable")
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float32To64 error %#v", err)
		}
	}()
	//用%v兼容f本身就是一个溢出的数值，可以打印出溢出部分
	s := fmt.Sprintf("%v", f)
	v, b := new(big.Float).SetString(s)
	if !b {
		return 0, false
	}
	f64, _ := v.Float64()
	return f64, true
}

// 计算两个浮点数相加
//
// float64 f1+f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func FloatAdd(f1, f2 float64) (*big.Float, bool) {
	//panic(ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("FloatAdd error %#v", err)
		}
	}()
	v := new(big.Float).Add(big.NewFloat(f1), big.NewFloat(f2))
	return v, true
}

// 计算两个浮点数相减
//
// float64 f1-f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func FloatSub(f1, f2 float64) (*big.Float, bool) {
	//panic(ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("FloatSub error %#v", err)
		}
	}()
	v := new(big.Float).Sub(big.NewFloat(f1), big.NewFloat(f2))
	return v, true
}

// 计算两个浮点数相乘
//
// float64 f1*f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func FloatMul(f1, f2 float64) (*big.Float, bool) {
	//panic(ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("FloatMul error %#v", err)
		}
	}()
	v := new(big.Float).Mul(big.NewFloat(f1), big.NewFloat(f2))
	return v, true
}

// 计算两个浮点数相除
//
// float64 f1/f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func FloatDiv(f1, f2 float64) (*big.Float, bool) {
	//panic(ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("FloatDiv error %#v", err)
		}
	}()
	v := new(big.Float).Quo(big.NewFloat(f1), big.NewFloat(f2))
	return v, true
}

// 计算两个浮点数相加
//
// float64 f1+f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float64Add(f1, f2 float64) (float64, bool) {
	//panic("unreachable")
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float64Add error %#v", err)
		}
	}()
	v, ok := FloatAdd(f1, f2)
	if !ok {
		return 0, false
	}
	f, _ := v.Float64()
	return f, true
}

// 计算两个浮点数相减
//
// float64 f1-f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float64Sub(f1, f2 float64) (float64, bool) {
	//panic("unreachable")
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float64Sub error %#v", err)
		}
	}()
	v, ok := FloatSub(f1, f2)
	if !ok {
		return 0, false
	}
	f, _ := v.Float64()
	return f, true
}

// 计算两个浮点数相乘
//
// float64 f1*f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float64Mul(f1, f2 float64) (float64, bool) {
	//ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float64Mul error %#v", err)
		}
	}()
	v, ok := FloatMul(f1, f2)
	if !ok {
		return 0, false
	}
	f, _ := v.Float64()
	return f, true
}

// 计算两个浮点数相除
//
// float64 f1/f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float64Div(f1, f2 float64) (float64, bool) {
	//ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float64Div error %#v", err)
		}
	}()
	v, ok := FloatDiv(f1, f2)
	if !ok {
		return 0, false
	}
	f, _ := v.Float64()
	return f, true
}

// 计算两个浮点数相加
//
// float32 f1+f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float32Add(f1, f2 float32) (float32, bool) {
	//panic("unreachable")
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float32Add error %#v", err)
		}
	}()
	f3, ok := Float32To64(f1)
	if !ok {
		return 0, false
	}
	f4, ok := Float32To64(f2)
	if !ok {
		return 0, false
	}
	v, ok := FloatAdd(f3, f4)
	if !ok {
		return 0, false
	}
	f, _ := v.Float32()
	return f, true
}

// 计算两个浮点数相减
//
// float32 f1-f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float32Sub(f1, f2 float32) (float32, bool) {
	//panic("unreachable")
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float32Sub error %#v", err)
		}
	}()
	f3, ok := Float32To64(f1)
	if !ok {
		return 0, false
	}
	f4, ok := Float32To64(f2)
	if !ok {
		return 0, false
	}
	v, ok := FloatSub(f3, f4)
	if !ok {
		return 0, false
	}
	f, _ := v.Float32()
	return f, true
}

// 计算两个浮点数相乘
//
// float32 f1*f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float32Mul(f1, f2 float32) (float32, bool) {
	//ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float32Mul error %#v", err)
		}
	}()
	f3, ok := Float32To64(f1)
	if !ok {
		return 0, false
	}
	f4, ok := Float32To64(f2)
	if !ok {
		return 0, false
	}
	v, ok := FloatMul(f3, f4)
	if !ok {
		return 0, false
	}
	f, _ := v.Float32()
	return f, true
}

// 计算两个浮点数相除
//
// float32 f1/f2结果
//
// bool 数值或计算异常时，返回false，计算正常时返回true
func Float32Div(f1, f2 float32) (float32, bool) {
	//ErrNaN
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Float32Div error %#v", err)
		}
	}()
	f3, ok := Float32To64(f1)
	if !ok {
		return 0, false
	}
	f4, ok := Float32To64(f2)
	if !ok {
		return 0, false
	}
	v, ok := FloatDiv(f3, f4)
	if !ok {
		return 0, false
	}
	f, _ := v.Float32()
	return f, true
}
