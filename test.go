package main

import (
	"errors"
	"fmt"
	"github.com/c4milo/unpackit"
	"os"
	"path"
	"time"
)

type Str struct {
	s string
}

type IInt int

func getNil() *Str {
	return nil
}
func getReturn() (str Str, err error) {
	return
}
func getReturn2() (str Str, err error) {
	err = errors.New("getReturn2 error")
	return
}
func getReturn3() (str []*Str, err error) {
	err = errors.New("func getReturn3() ")
	fmt.Println(&err)
	return
}
func getReturn4() (str []*Str, err *error) {
	err2 := errors.New("func getReturn3() ")
	fmt.Println(&err2)
	return nil, &err2
}
func getInt() *int {
	return nil
}
func getInt2() *IInt {
	return nil
}
func getInt3() (is []*IInt) {
	fmt.Println(len(is))
	return nil
}
func test3() (data map[int]*Str, err error) {
	fmt.Println(data, err)
	return
}

type Struct1 struct {
	test string
}

func (s *Struct1) name() {
	fmt.Println(s.test)
}

type Struct2 struct {
	Struct1
}

func (s *Struct2) name() {
	fmt.Println("---------", s.test)
}
func main() {
	fmt.Println(path.Base("https://gorm.io/zh_CN/docs/sql_builder.html"))

	var name map[int]*Str
	if arr, ok := name[2]; ok {
		fmt.Println("3223")
		fmt.Println("2333", arr)
	}
	getInt3()
	data, err := test3()
	fmt.Println(data, err)

	{
		d := -1 * time.Duration(847) * time.Second
		fmt.Println(d, "小时")
	}
	{
		d := -1 * time.Duration(84585) * time.Second
		fmt.Println(d, "小时")
	}
}

func getReturn4Test() {
	str, err := getReturn4()
	fmt.Println(str, err, *err)
}

func funcName() {
	s, err := getReturn()
	fmt.Println("getReturn s:", s.s)
	if nil != err {
		fmt.Println("getReturn err:", err)
	}
	s2, err2 := getReturn2()
	fmt.Println("getReturn s2:", s2.s)
	if nil != err2 {
		fmt.Println("getReturn err:", err2)
	}
	str := getNil()
	if str.s == "" {
		fmt.Println("no nil pointer error")
		return
	}
	fmt.Println("unknown error")
}

func unzip() {
	file, _ := os.Open("data/log/ocean_weixin_push_group_5_2024-04-27.txt.gz")
	_ = unpackit.Unpack(file, "/path/to/file")
}
