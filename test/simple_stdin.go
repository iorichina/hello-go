package test

import (
	"errors"
	"fmt"
	"os"
)

func stdin(min int) (in string, err error) {
	//reader := bufio.NewReader(os.Stdin)
	//result, err := reader.ReadString('\n')
	//if err != nil {
	//	fmt.Println("read error:", err)
	//}

	fmt.Println("请输入：")
	var buffer [512]byte
	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		fmt.Println("输入错误：", err)
		return "", err
	}
	if n < min {
		fmt.Println("输入字符少于", min)
		return "", errors.New("输入字符少于" + string(rune(min)))
	}
	return string(buffer[:]), nil
}
