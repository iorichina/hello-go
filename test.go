package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c4milo/unpackit"
	"net"
	"os"
	"strings"
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

func setStruct1Name1(s *Struct1, name string) {
	s.test = name
}

func setStruct1Name2(s *Struct1, name string) {
	(*s).test = name
}

type Struct2 struct {
	Struct1
}

func (s *Struct2) name() {
	fmt.Println("---------", s.test)
}

var e = errors.New("")

func throwE() error {
	return e
}

type QrCodeStatus int //二维码状态

const (
	qrCodeStatusInit  QrCodeStatus = 0 //订单初始化
	qrCodeStatusScan  QrCodeStatus = 1 //已扫码
	qrCodeStatusPayed QrCodeStatus = 3 //已支付
)

func printQrCode(status QrCodeStatus) {
	fmt.Println(status)
	switch status {
	case qrCodeStatusInit:
		fmt.Println(1110)

	case qrCodeStatusScan:
		fmt.Println(1111)

	case qrCodeStatusPayed:
		fmt.Println(1113)
	default:
		fmt.Println(1119)
	}
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}

func getBytes() ([]byte, error) {
	return nil, nil
}

func main() {
	bs := []byte{0x38, 0x34, 0x43, 0x32, 0x45, 0x34, 0x46, 0x36, 0x42, 0x42, 0x45, 0x41}
	fmt.Println(string(bs))
	v := &net.OpError{}
	marshal, _ := json.Marshal(v)
	fmt.Printf("%#v\n", v)
	fmt.Printf("%#v\n", marshal)
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
