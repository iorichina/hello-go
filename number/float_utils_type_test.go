package number

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestFloat32Type64(t *testing.T) {
	{
		fmt.Println(1)
		var f float32 = 3.1415926
		fmt.Printf("%f\t%v\t%T\n", f, f, f)
		fmt.Println(f)
		fmt.Println(fmt.Sprint(f))
	}
	{
		fmt.Println(2)
		var f float32 = 3.1415926
		f1 := float64(f)
		fmt.Printf("%f\t%v\t%T\n", f1, f1, f1)
		fmt.Println(f1)
		fmt.Println(fmt.Sprint(f1))
	}
	{
		fmt.Println(3)
		var f float32 = 3.141599
		fmt.Printf("%f\t%v\t%T\n", f, f, f)
		fmt.Println(f)
		fmt.Println(fmt.Sprint(f))
	}
	{
		fmt.Println(4)
		var f float32 = 3.141599
		f1 := float64(f)
		fmt.Printf("%f\t%v\t%T\n", f1, f1, f1)
		fmt.Println(f1)
		fmt.Println(fmt.Sprint(f1))
	}
	{
		fmt.Println(5)
		var f float32 = 3.141592502593994
		f1 := float64(f)
		fmt.Printf("%f\t%v\t%T\n", f1, f1, f1)
		fmt.Println(f1)
		fmt.Println(fmt.Sprint(f1))
	}
	{
		fmt.Println(6)
		var f float32 = 1.123456789
		fmt.Printf("%f\t%v\t%T\n", f, f, f)
		fmt.Println(f)
		fmt.Println(fmt.Sprint(f))
	}

	fmt.Printf("%f\t%v\n", math.MaxFloat32, math.MaxFloat32)
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
