package test

import (
	"encoding/base64"
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
		fmt.Println("输入字符太少")
		return
	}
	s := string(buffer[:n])
	res, err := decodeWebStr(s)
	fmt.Println("解析结果：")
	fmt.Println(res)
	fmt.Println()
	fmt.Print("按回车退出...")
	os.Stdin.Read(buffer[:])
}

func decodeWebStr(enc3 string) (string, error) {
	encDescDefaultWeb := "dX9f9#q_9XG"
	encString, err := fromEncString(enc3, []byte(encDescDefaultWeb))
	if nil != err {
		fmt.Println("fromEncString报错:", encString, err)
		return "", err
	}
	return string(encString), nil
}

func fromEncString(enc3 string, encDesc []byte) ([]byte, error) {
	if encDesc == nil {
		return []byte(enc3), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(enc3)
	if nil != err {
		fmt.Println("base64解析报错:", enc3, err)
		return nil, err
	}
	enc2 := string(decoded)
	index := len([]rune(enc2)) - 1
	if index > 30 {
		index = 30
	}
	enc1 := substring(enc2, 0, index) + substring(enc2, index+1, -1)
	res, err := base64.StdEncoding.DecodeString(enc1)
	if nil != err {
		fmt.Println("除灵后的base64解析报错:", enc1, err)
		return nil, err
	}
	length := len(res)
	lenDesc := len(encDesc)
	for i := 0; i < length; i++ {
		res[i] ^= encDesc[i%lenDesc]
	}
	return res, nil
}

func substring(str string, beginInclude int, endExclude int) string {
	if beginInclude < 0 {
		beginInclude = 0
	}
	if -1 != endExclude && endExclude <= beginInclude {
		return ""
	}
	rs := []rune(str)
	length := len(rs)
	if beginInclude >= length {
		return ""
	}
	if -1 == endExclude {
		endExclude = length
	}
	if endExclude > length {
		endExclude = length
	}
	return string(rs[beginInclude:endExclude])
}
