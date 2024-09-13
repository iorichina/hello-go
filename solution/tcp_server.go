package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// go build tcp_server.go
//
// [nohup] sudo ./tcp_server >>log.log 2>&1 &
//
// TCP Server端测试
// 处理函数
func process(logger *log.Logger, conn net.Conn) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	var mac string

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			now := time.Now().Format("2006-01-02 15:04:05.000")
			logger.Printf("%v[%v][%v]server Close err %v\t%#v\n", now, conn.RemoteAddr(), mac, err.Error(), err)
		}
	}(conn) // 关闭连接
	now = time.Now().Format("2006-01-02 15:04:05.000")
	remote := conn.RemoteAddr()
	logger.Printf("%v[%v]Connected to %v\n", now, remote, conn.LocalAddr())

	var err error
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 11, 48, 90, 0, 149})
	if err != nil {
		logger.Printf("%v[%v]server Write(0x30) err %v\n", now, remote, err)
	}
	now = time.Now().Format("2006-01-02 15:04:05.000")
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 51, 60})
	if err != nil {
		logger.Printf("%v[%v]server Write(0x33) err %v\n", now, remote, err)
	}
	now = time.Now().Format("2006-01-02 15:04:05.000")
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 66, 75})
	if err != nil {
		logger.Printf("%v[%v]server Write(0x42) err %v\n", now, remote, err)
	}
	now = time.Now().Format("2006-01-02 15:04:05.000")
	_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
	if err != nil {
		logger.Printf("%v[%v]server Write(0x34) err %v\n", now, remote, err)
	}

	var n int
	reader := bufio.NewReader(conn)
	var buf [512]byte
	for {
		n, err = reader.Read(buf[:]) // 读取数据
		if err != nil {
			now = time.Now().Format("2006-01-02 15:04:05.000")
			if errors.Is(err, io.EOF) {
				logger.Printf("%v[%v][%v]server Read EOF %#v\n", now, remote, mac, err)
			} else {
				logger.Printf("%v[%v][%v]server Read failed err %v\t%#v\n", now, remote, mac, err.Error(), err)
			}
			break
		}

		now = time.Now().Format("2006-01-02 15:04:05.000")
		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x34 != buf[7] && 0x35 != buf[7] {
				continue
			}
			if 0x34 == buf[7] {
				mac = strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
				logger.Printf("%v[%v][%v]server 读 %#x 状态 %d (0=idle,1=playing,>1=error)\n", now, remote, mac, buf[7], buf[8])
				continue
			}
			if 0x35 == buf[7] {
				mac = strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
				logger.Printf("%v[%v][%v]server 读 %#x 心跳\n", now, remote, mac, buf[7])

				_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 51, 60})
				if err != nil {
					logger.Printf("%v[%v][%v]server Write(0x33) err %v\n", now, remote, mac, err)
				}
				_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 11, 49, 0, 1, 61})
				if err != nil {
					logger.Printf("%v[%v][%v]server Write(0x31) err %v\n", now, remote, mac, err)
					continue
				}
				_, err = conn.Write([]byte{254, 73, 66, 1, 182, 189, 11, 49, 0, 1, 61})
				if err != nil {
					logger.Printf("%v[%v][%v]server Write(0x30) err %v\n", now, remote, mac, err)
					continue
				}
				_, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 10, 50, 2, 62})
				if err != nil {
					logger.Printf("%v[%v][%v]server Write(0x32) err %v\n", now, remote, mac, err)
					continue
				}
				continue
			}
			logger.Printf("%v[%v][%v]server 读 %#x\t% X\n", now, remote, mac, buf[7], buf[:n])
			continue
		}

		recvStr := string(buf[:n])
		logger.Printf("%v[%v]server 读 %v\t% X\n", now, remote, recvStr, buf[:n])
		_, err = conn.Write([]byte(recvStr))
		if err != nil {
			logger.Printf("%v[%v][%v]server Reply err %v\n", now, remote, mac, err)
			continue
		} // 发送数据
	}
}

func main() {
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatalf("create file server.log failed: %v", err)
	}
	logger := log.New(io.MultiWriter(os.Stdout, f), "", 0)
	listen, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		logger.Printf("%v Listen() failed, err %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)
		return
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			logger.Printf("%v Accept() failed, err: %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)
			continue
		}
		go process(logger, conn) // 启动一个goroutine来处理客户端的连接请求
	}
}
