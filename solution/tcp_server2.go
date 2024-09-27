package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// go build tcp_server2.go
//
// vi run-server2.sh
//
// ./tcp_server >>server2.log 2>&1 &
//
// sudo sh run-server2.sh
//
// TCP Server端测试
// 处理函数
func main() {
	logger := log.New(os.Stdout, "", log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
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
		go process2(conn) // 启动一个goroutine来处理客户端的连接请求
	}
}

func process2(conn net.Conn) {
	remote := conn.RemoteAddr()
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v][%v]", mac, remote), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	defer func() {
		if r := recover(); r != nil {
			logger.Printf("Recovered from panic: %#v", r)
		}
	}()

	defer func() {
		err = conn.Close()
		if err != nil {
			logger.Printf("Connection Close err %v\n", err)
			return
		}
		logger.Printf("Connection Close\n")
	}() // 关闭连接
	logger.Printf("Connection Connected to %v\n", conn.LocalAddr())

	queueChan := make(chan []byte, 1024)
	defer close(queueChan)
	go func() {
		for bb := range queueChan {
			_, err = conn.Write(bb)
			if len(bb) > 8 && 0xFE == bb[0] && 0x01 == bb[3] {
				if 0x31 == bb[7] || 0x34 == bb[7] || 0x35 == bb[7] {
					logger.Printf("Write %#v with %v\n", bb[7], err)
				}
			}
		}
		queueChan = nil
	}()

	heartbeatChan := make(chan []byte)
	defer close(heartbeatChan)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("Recovered heartbeat from panic: %v\n", r)
			}
		}()
		for {
			if nil == queueChan {
				break
			}
			select {
			case <-heartbeatChan:
			case <-time.After(15 * time.Second):
			}
			if nil == queueChan {
				break
			}
			queueChan <- []byte{254, 134, 226, 1, 121, 29, 9, 0x33, 60}
			queueChan <- []byte{254, 73, 66, 1, 182, 189, 11, 0x31, 0, 1, 61}
			queueChan <- []byte{254, 134, 226, 1, 121, 29, 10, 0x32, 2, 62}
			queueChan <- []byte{254, 134, 226, 1, 121, 29, 9, 0x34, 61}
		}
	}()

	queueChan <- []byte{254, 134, 226, 1, 121, 29, 11, 0x30, 90, 0, 149}
	queueChan <- []byte{254, 134, 226, 1, 121, 29, 9, 0x33, 60}
	queueChan <- []byte{254, 134, 226, 1, 121, 29, 9, 0x42, 75}
	queueChan <- []byte{254, 134, 226, 1, 121, 29, 9, 0x34, 61}

	scanner := newServerScanner(conn)
	for {
		if !scanner.Scan() {
			err = scanner.Err()
			if err != nil {
				logger.Printf("Connection Read failed %v\n", err)
			} else {
				err = io.EOF
				logger.Printf("Connection Read EOF %v\n", err)
			}
			break
		}
		buf := scanner.Bytes()
		if len(buf) <= 0 {
			continue
		}

		if len(buf) > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x31 == buf[7] {
				logger.Printf("Read %#x\t% X\n", buf[7], buf[8:10])
				continue
			}
			if 0x34 == buf[7] {
				m := strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
				if m != mac {
					logger = log.New(os.Stdout, fmt.Sprintf("[%17v][%v]", m, remote), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
				}
				mac = m
				logger.Printf("Read %#x with %d (0=idle,1=playing,>1=error)\n", buf[7], buf[8])
				if buf[8] == 0 {
					queueChan <- []byte{254, 134, 226, 1, 121, 29, 11, 0x30, 90, 0, 149}
				}
				continue
			}
			if 0x35 == buf[7] {
				m := strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
				if m != mac {
					logger = log.New(os.Stdout, fmt.Sprintf("[%17v][%v]", m, remote), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
				}
				mac = m
				logger.Printf("Read %#x\n", buf[7])
				heartbeatChan <- nil
				continue
			}
			continue
		}

		recvStr := string(buf[:])
		logger.Printf("读 %v\t% X\n", recvStr, buf[:])
	}
}

// 命令头	消息ID高位	消息ID低位	命令头取反	消息ID高位取反	消息ID低位取反	包长度	指令码	数据	校验位
// 0xfe		0x00	   0x01		  0x01		 0xff		    0xfe			0x0a   0x14	  Data	  sum(包长度+指令码+数据...)%256
func newServerScanner(rd io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(rd)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if len(data) > 8 && 0xFE == data[0] && 0x01 == data[3] {
			length := int(data[6])
			return length, data[:length], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	return scanner
}
