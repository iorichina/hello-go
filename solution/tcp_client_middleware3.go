package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) <= 4 {
		log.Fatalf("need local ip-port and remote ip-port")
	}
	logger := log.New(os.Stdout, "main   ", log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	localAddr := os.Args[1]
	localTimeout, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		logger.Fatalf("parse local timeout fail,err:%s", err)
	}
	remoteAddr := os.Args[3]
	remoteTimeout, err := strconv.ParseInt(os.Args[4], 10, 64)
	if err != nil {
		logger.Fatalf("parse remote timeout fail,err:%s", err)
	}

	localTimeoutDuration := time.Duration(localTimeout) * time.Millisecond
	localConn, err := net.DialTimeout("tcp", localAddr, localTimeoutDuration)
	if nil != err {
		log.Fatalf("local ip-port %v unreachable err %v", localAddr, err)
	}
	logger.Printf("[%v]local Connection Connected\n", localAddr)

	remoteTimeoutDuration := time.Duration(remoteTimeout) * time.Millisecond
	remoteConn, err := net.DialTimeout("tcp", remoteAddr, remoteTimeoutDuration)
	if nil != err {
		log.Fatalf("remote ip-port %v unreachable err %v", remoteAddr, err)
	}
	logger.Printf("[%v]remote Connection Connected\n", remoteAddr)

	clientMiddleware3(localAddr, localTimeoutDuration, remoteAddr, remoteTimeoutDuration, localConn, remoteConn)
	logger.Printf("([%v]->[%v])process stoped\n", localAddr, remoteAddr)
}

func clientMiddleware3(localAddr string, localTimeoutDuration time.Duration, remoteAddr string, remoteTimeoutDuration time.Duration, localConn, remoteConn net.Conn) {
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v]([%v]->[%v])middle ", mac, localAddr, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	defer func(conn net.Conn) {
		if nil == conn {
			return
		}
		err := conn.Close()
		logger.Printf("Connection local Close %v\n", err)
	}(localConn) // 关闭连接
	defer func(conn net.Conn) {
		if nil == conn {
			return
		}
		err := conn.Close()
		logger.Printf("Connection remote Close %v\n", err)
	}(remoteConn) // 关闭连接

	localChan := make(chan error)
	defer close(localChan)
	remoteChan := make(chan error)
	defer close(remoteChan)

	localQueue := make(chan []byte, 1024)
	remoteQueue := make(chan []byte, 1024)

	macChan := make(chan string, 128)
	macChanLocal := make(chan string, 128)
	macChanRemote := make(chan string, 128)
	go func() {
		for m := range macChan {
			if m != mac {
				logger = log.New(os.Stdout, fmt.Sprintf("[%17v]([%v]->[%v])middle ", m, localAddr, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
			}
			mac = m
			macChanLocal <- m
			macChanRemote <- m
		}
		defer close(macChanLocal)
		defer close(macChanRemote)
	}()

	go handleLocal3(localAddr, remoteAddr, macChan, macChanLocal, localConn, localQueue, remoteQueue, localChan, remoteChan)
	go handleRemote3(localAddr, remoteAddr, macChan, macChanRemote, remoteConn, localQueue, remoteQueue, localChan, remoteChan)

	select {
	case err = <-localChan:
		logger.Printf("Connection local err, %v, stop process\n", err)
		return
	case err = <-remoteChan:
		logger.Printf("Connection remote err, %v, stop process\n", err)
		return
	}
}

func handleLocal3(localAddr, remoteAddr string, macChan, macChanLocal chan string, localConn net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	defer close(remoteQueue)
	defer close(macChan)
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v][%v]local  ", mac, localAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	go func() {
		for m := range macChanLocal {
			if m != mac {
				logger = log.New(os.Stdout, fmt.Sprintf("[%17v][%v]local  ", m, localAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
			}
			mac = m
		}
	}()

	go func() {
		for bb := range localQueue {
			conn := localConn
			if nil == conn {
				if len(bb) > 8 && 0xFE == bb[0] && 0x01 == bb[3] {
					if 0x31 == bb[7] {
						logger.Printf("Fallback %#v\n", bb[7])
						bb[10] = bb[10] - bb[9] - bb[8]
						bb[8] = 0x00
						bb[9] = 0x00
						remoteQueue <- bb
					}
					//todo 0x34
				}
				continue
			}
			_, err = conn.Write(bb)
			if err != nil {
			}
			if len(bb) > 8 && 0xFE == bb[0] && 0x01 == bb[3] {
				if 0x31 == bb[7] || 0x34 == bb[7] || 0x35 == bb[7] {
					logger.Printf("Write %#v with %v\n", bb[7], err)
				}
			}
		}
	}()

	_, err = localConn.Write(newClientMiddleware3Msg(0x34, nil))
	logger.Printf("Write(0x34) connect with %v\n", err)
	scanner := newClientMiddleware3Scanner(localConn)
	for {
		if !scanner.Scan() {
			err = scanner.Err()
			if err == nil {
				err = io.EOF
			}
			localChan <- err
			break
		}
		buf := scanner.Bytes()
		if len(buf) <= 0 {
			continue
		}

		if len(buf) > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x34 == buf[7] {
				m := strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
				if m != mac {
					logger = log.New(os.Stdout, fmt.Sprintf("[%17v][%v]local  ", m, localAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
				}
				mac = m
				macChan <- m
				logger.Printf("Read %#v(%d) with %d (0=idle,1=playing,>1=error)\n", buf[7], int(buf[1])*256+int(buf[2]), buf[8])
			} else if 0x35 == buf[7] {
				m := strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
				if m != mac {
					logger = log.New(os.Stdout, fmt.Sprintf("[%17v][%v]local  ", m, localAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
				}
				mac = m
				macChan <- m
				logger.Printf("Read %#v(%d)\n", buf[7], int(buf[1])*256+int(buf[2]))
			} else if 0x31 == buf[7] {
				logger.Printf("Read %#v(%d) with % X\n", buf[7], int(buf[1])*256+int(buf[2]), buf[8:10])
			}
		}

		dup := make([]byte, len(buf))
		copy(dup, buf[:])
		remoteQueue <- dup
	}
}

func handleRemote3(localAddr, remoteAddr string, macChan, macChanRemote chan string, remoteConn net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	defer close(localQueue)
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v][%v]remote ", mac, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	go func() {
		for m := range macChanRemote {
			if m != mac {
				logger.SetPrefix(fmt.Sprintf("[%17v][%v]remote ", m, remoteAddr))
			}
			mac = m
		}
	}()

	go func() {
		for buf := range remoteQueue {
			conn := remoteConn
			if nil == conn {
				continue
			}
			_, err = conn.Write(buf)
			if err != nil {
			} // 发送数据
			if len(buf) > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
				if 0x31 == buf[7] || 0x34 == buf[7] || 0x35 == buf[7] {
					logger.Printf("Write %#x(%d) to local with %v\n", buf[7], int(buf[1])*256+int(buf[2]), err)
				}
			}
		}
	}()

	scanner := newClientMiddleware3Scanner(remoteConn)
	for {
		if !scanner.Scan() {
			err = scanner.Err()
			if err == nil {
				err = io.EOF
			}
			remoteChan <- err
			continue
		}
		buf := scanner.Bytes()
		if len(buf) <= 0 {
			continue
		}

		if len(buf) > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x31 == buf[7] || 0x34 == buf[7] || 0x35 == buf[7] {
				logger.Printf("Read %#v(%d)\n", buf[7], int(buf[1])*256+int(buf[2]))
			}
		}

		dup := make([]byte, len(buf))
		copy(dup, buf[:])
		localQueue <- dup
	}
}

// 命令头	消息ID高位	消息ID低位	命令头取反	消息ID高位取反	消息ID低位取反	包长度	指令码	数据	校验位
// 0xfe		0x00	   0x01		  0x01		 0xff		    0xfe			0x0a   0x14	  Data	  sum(包长度+指令码+数据...)%256
func newClientMiddleware3Scanner(rd io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(rd)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if len(data) > 8 && 0xFE == data[0] && 0x01 == data[3] {
			length := int(data[6])
			if length >= len(data) {
				return length, data[:length], nil
			}
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	return scanner
}

func newClientMiddleware3Msg(cmd byte, data []byte) []byte {
	length := 6 + 1 + 1 + len(data) + 1
	msg := make([]byte, length)
	id := rand.Int() & 0xFFFF
	msg[0] = 0xFE
	msg[1] = byte(id >> 8)
	msg[2] = byte(id & 0xFF)
	msg[3] = 0x01
	msg[4] = ^msg[1]
	msg[5] = ^msg[2]
	msg[6] = byte(length)
	msg[7] = cmd
	sum := int(msg[6]) + int(msg[7])
	if nil != data && len(data) > 0 {
		for i, v := range data {
			msg[8+i] = v
			sum += int(v)
		}
	}
	msg[length-1] = byte(sum & 0xFF)
	return msg
}
