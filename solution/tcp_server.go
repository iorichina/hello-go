package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// go build tcp_server.go
//
// vi run-server.sh
//
// ./tcp_server >>server.log 2>&1 &
//
// sudo sh run-server.sh
//
// TCP Server端测试
// 处理函数
func process(logger *log.Logger, conn net.Conn) {
	remote := conn.RemoteAddr()
	var mac string

	defer func() {
		if r := recover(); r != nil {
			logger.Printf("[%v][%v]Recovered from panic: %#v", remote, mac, r)
		}
	}()

	defer func() {
		err := conn.Close()
		logger.Printf("[%v][%v]Connection Close\n", remote, mac)
		if err != nil {
			logger.Printf("[%v][%v]Close err %v\n", remote, mac, err.Error())
		}
	}() // 关闭连接
	logger.Printf("[%v][%v]Connection Connected to %v\n", remote, mac, conn.LocalAddr())

	var err error
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 11, 48, 90, 0, 149})
	if err != nil {
		logger.Printf("[%v]Write(0x30) err %v\n", remote, err)
	}
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 51, 60})
	if err != nil {
		logger.Printf("[%v]Write(0x33) err %v\n", remote, err)
	}
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 66, 75})
	if err != nil {
		logger.Printf("[%v]Write(0x42) err %v\n", remote, err)
	}
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
	if err != nil {
		logger.Printf("[%v]Write(0x34) err %v\n", remote, err)
	}

	var n int
	reader := bufio.NewReader(conn)
	var buf [512]byte
	for {
		n, err = reader.Read(buf[:]) // 读取数据
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Printf("[%v][%v]Read EOF %#v\n", remote, mac, err)
			} else {
				logger.Printf("[%v][%v]Read failed err %v\n", remote, mac, err.Error())
			}
			break
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x31 != buf[7] && 0x34 != buf[7] && 0x35 != buf[7] {
				continue
			}
			if 0x31 == buf[7] {
				logger.Printf("[%v][%v]读 %#x\t% X\n", remote, mac, buf[7], buf[8:10])
				continue
			}
			if 0x34 == buf[7] {
				mac = strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
				logger.Printf("[%v][%v]读 %#x 状态 %d (0=idle,1=playing,>1=error)\n", remote, mac, buf[7], buf[8])
				continue
			}
			if 0x35 == buf[7] {
				mac = strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")

				_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 51, 60})
				if err != nil {
					logger.Printf("[%v][%v]Write(0x33) err %v\n", remote, mac, err)
				}
				_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 11, 49, 0, 1, 61})
				if err != nil {
					logger.Printf("[%v][%v]Write(0x31) err %v\n", remote, mac, err)
					continue
				}
				_, err = conn.Write([]byte{254, 73, 66, 1, 182, 189, 11, 49, 0, 1, 61})
				if err != nil {
					logger.Printf("[%v][%v]Write(0x31) err %v\n", remote, mac, err)
					continue
				}
				_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 10, 50, 2, 62})
				if err != nil {
					logger.Printf("[%v][%v]Write(0x32) err %v\n", remote, mac, err)
					continue
				}
				continue
			}
			logger.Printf("[%v][%v]读 %#x\t% X\n", remote, mac, buf[7], buf[:n])
			continue
		}

		recvStr := string(buf[:n])
		logger.Printf("[%v]读 %v\t% X\n", remote, recvStr, buf[:n])
		// _, err = conn.Write([]byte(recvStr))
		// if err != nil {
		// 	logger.Printf("[%v][%v]Reply err %v\n", remote, mac, err)
		// 	continue
		// } // 发送数据
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	listen, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		logger.Printf("Listen() failed, err %#v\n", err)
		return
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			logger.Printf("Accept() failed, err: %#v\n", err)
			continue
		}
		go process(logger, conn) // 启动一个goroutine来处理客户端的连接请求
	}
}
