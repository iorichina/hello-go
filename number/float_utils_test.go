package number

import (
	"fmt"
	"hello-go/logger"
	"strconv"
	"testing"
)

func TestStrToFloat64(t *testing.T) {
	s := "3.14159"
	f, ok := StrToFloat64(s)
	if !ok {
		t.Fatalf("TestStrToFloat64 失败\n")
	}
	if 3.14159 != f {
		t.Fatalf("TestStrToFloat64 失败 %#v\n", f)
	}
}

func TestFloat32To64(t *testing.T) {
	var f float32 = 3.1415926 //本身就是溢出
	fmt.Printf("%f\t%v\t%T\n", f, f, f)

	// float32 转 float64
	fmt.Printf("%v\n", float64(f)) // 输出：3.1415925

	// float64 转 float32
	var f2 float64 = 3.141592653589793
	fmt.Printf("%v\n", float32(f2)) // 输出：3.1415927

	to64, b := Float32To64(f)
	if !b {
		t.Fatal("Float32To64 失败")
	}
	var f64 float64 = 3.1415926
	fmt.Printf("to64 %f\t%v\t%T\n", to64, to64, to64)
	fmt.Printf("f64 %f\t%v\t%T\n", f64, f64, f64)
	if to64 == f64 {
		t.Fatal("Float32To64 失败")
	}
	float, err := strconv.ParseFloat(fmt.Sprint(f), 64)
	fmt.Printf("ParseFloat %f\t%v\t%T\n", float, float, float)
	if err != nil {
		t.Fatal("ParseFloat 失败")
	}
	if to64 != float {
		t.Fatal("Float32To64 失败")
	}
}

func TestFloat64Utils(t *testing.T) {
	{
		add, ok := Float64Add(1, 2)
		if !ok {
			t.Fatal("Float64Add 失败")
		}
		if add != 3 {
			t.Fatal("Float64Add 失败")
		}
	}
	{
		add, ok := Float64Add(0.99, 0.01)
		if !ok {
			t.Fatal("Float64Add 失败")
		}
		if add != 1 {
			t.Fatal("Float64Add 失败")
		}
	}
}

func TestFloat32Utils(t *testing.T) {
	{
		var a float32
		a = 990
		b := a / 100 * 100
		if b == 990 {
			t.Fatal("预期失败案例执行失败")
		}
	}
	{
		var a float64
		a = 990
		b := a / 100 * 100
		if b != 990 {
			t.Fatal("预期失败案例执行失败")
		}
	}
	{
		var a float32
		a = 990
		div, b := Float32Div(a, 100)
		logger.Infof("Float32Div(%#v, %#v) %#v %#v", a, 100, div, b)
		if !b {
			t.Fatal("Float32Div 失败")
		}
		mul, b2 := Float32Mul(div, 100)
		logger.Infof("Float32Mul(%#v, %#v) %#v %#v", div, 100, mul, b2)
		if !b2 {
			t.Fatal("Float32Mul 失败")
		}
		if mul != 990 {
			t.Fatal("Float32Utils 失败")
		}
	}
	{
		div, b := FloatDiv(0, 0)
		fmt.Println(div, b)
		if b {
			t.Fatal("FloatDiv 失败")
		}
	}
}
