package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// TCP 客户端
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		fmt.Printf("%v Dial() failed, err %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%v[%v] Close() err %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), conn.RemoteAddr(), err)
		}
	}(conn) // 关闭TCP连接
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			return
		}

		var ss []byte
		//fe 86 e2 01 79 1d 0b 31 00 01 3d
		//fe 49 42 01 b6 bd 0b 31 00 01 3d
		//fe 86 e2 01 79 1d 0a 32 02 3e
		if 'f' == input[0] && 'e' == input[1] && ' ' == input[2] && len(strings.Split(inputInfo, " ")) > 8 {
			ss = make([]byte, 0)
			split := strings.Split(inputInfo, " ")
			for _, s := range split {
				i, _ := strconv.ParseInt(s, 16, 16)
				ss = append(ss, byte(i))
			}
		} else {
			ss = []byte(inputInfo)
		}
		_, err := conn.Write(ss) // 发送数据
		if err != nil {
			fmt.Printf("%v[%v]Write(%v) err %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), conn.RemoteAddr(), inputInfo, err)
			return
		}

		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("%v[%v]Read() err %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), conn.RemoteAddr(), err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
