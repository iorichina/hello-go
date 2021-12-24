package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("请输入：")
	var buffer [512]byte
	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		fmt.Println("输入错误：", err)
		return
	}
	if n < 10 {
		fmt.Println("输入字符少于", 10)
		return
	}
	s := string(buffer[:])
	fmt.Println(s)
	fmt.Println()
	fmt.Print("按回车退出...")
	os.Stdin.Read(buffer[:])
}
