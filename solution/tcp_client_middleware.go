package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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
	defer func(conn net.Conn) {
		if nil == conn {
			return
		}
		_ = conn.Close()
	}(localConn) // 关闭连接
	logger.Printf("[%v]local Connection Connected\n", localAddr)

	remoteTimeoutDuration := time.Duration(remoteTimeout) * time.Millisecond
	remoteConn, err := net.DialTimeout("tcp", remoteAddr, remoteTimeoutDuration)
	if nil != err {
		log.Fatalf("remote ip-port %v unreachable err %v", remoteAddr, err)
	}
	defer func(conn net.Conn) {
		if nil == conn {
			return
		}
		_ = conn.Close()
	}(remoteConn) // 关闭连接
	logger.Printf("[%v]remote Connection Connected\n", remoteAddr)

	clientMiddleware(localAddr, localTimeoutDuration, remoteAddr, remoteTimeoutDuration, localConn, remoteConn)
	logger.Printf("([%v]->[%v])process stoped\n", localAddr, remoteAddr)
}

func clientMiddleware(localAddr string, localTimeoutDuration time.Duration, remoteAddr string, remoteTimeoutDuration time.Duration, localConn, remoteConn net.Conn) {
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v]([%v]->[%v])middle ", mac, localAddr, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	defer func(conn net.Conn) {
		if nil == conn {
			return
		}
		err := conn.Close()
		logger.Printf("Connection local defer Close %v\n", err)
	}(localConn) // 关闭连接
	defer func(conn net.Conn) {
		if nil == conn {
			return
		}
		err := conn.Close()
		logger.Printf("Connection remote defer Close %v\n", err)
	}(remoteConn) // 关闭连接

	localChan := make(chan error)
	remoteChan := make(chan error)

	localQueue := make(chan []byte, 1024)
	defer close(localQueue)
	remoteQueue := make(chan []byte, 1024)
	defer close(remoteQueue)

	macChan := make(chan string, 128)
	macChanLocal := make(chan string, 128)
	macChanRemote := make(chan string, 128)
	go func() {
		for m := range macChan {
			if m != mac {
				logger.SetPrefix(fmt.Sprintf("[%17v]([%v]->[%v])middle ", m, localAddr, remoteAddr))
			}
			mac = m
			macChanLocal <- m
			macChanRemote <- m
		}
		defer close(macChanLocal)
		defer close(macChanRemote)
	}()

	localConnChan := make(chan net.Conn)
	defer close(localConnChan)
	remoteConnChan := make(chan net.Conn)
	defer close(remoteConnChan)

	go handleLocal(localAddr, remoteAddr, macChan, macChanLocal, localConn, localConnChan, localQueue, remoteQueue, localChan, remoteChan)
	go handleRemote(localAddr, remoteAddr, macChan, macChanRemote, remoteConn, remoteConnChan, localQueue, remoteQueue, localChan, remoteChan)

	for {
		select {
		case err = <-localChan:
			start := time.Now()
			logger.Printf("Connection local err, %v, trying to reconnect\n", err)
			localConn, err = net.DialTimeout("tcp", localAddr, localTimeoutDuration)
			logger.Printf("Connection local reconnect %v with %v\n", localConn, err)
			duration := localTimeoutDuration - time.Since(start)
			if nil != err && duration > 0 {
				time.Sleep(duration)
			}
			localConnChan <- localConn
		case err = <-remoteChan:
			start := time.Now()
			logger.Printf("Connection remote err, %v, trying to reconnect\n", err)
			remoteConn, err = net.DialTimeout("tcp", remoteAddr, remoteTimeoutDuration)
			logger.Printf("Connection remote reconnect %v with %v\n", remoteConn, err)
			duration := remoteTimeoutDuration - time.Since(start)
			if nil != err && duration > 0 {
				time.Sleep(duration)
			}
			remoteConnChan <- remoteConn
		}
	}
}

func handleLocal(localAddr, remoteAddr string, macChan, macChanLocal chan string, localConn net.Conn, localConnChan chan net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	defer close(localChan)
	defer close(macChan)
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v][%v]local  ", mac, localAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	go func() {
		for m := range macChanLocal {
			if m != mac {
				logger.SetPrefix(fmt.Sprintf("[%17v][%v]local  ", m, localAddr))
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

	_, err = localConn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
	logger.Printf("Write(0x34) by connected with %v\n", err)
	reader := bufio.NewReader(localConn)
	var buf [1024]byte
	var n int
	for {
		if nil == localConn {
			for {
				if nil == localConn {
					localConn = <-localConnChan
				}
				if nil == localConn {
					localChan <- errors.New("need retry")
					continue
				}
				// reader.Reset(localConn)
				reader = bufio.NewReader(localConn)
				_, err = localConn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
				logger.Printf("Write(0x34) reconnect with %v\n", err)
				break
			}
		}
		if n <= 0 {
			continue
		}

		n, err = reader.Read(buf[:]) // todo 粘包
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Printf("Connection Read EOF %v\n", err)
			} else {
				logger.Printf("Connection Read failed %v\n", err)
			}
			conn := localConn
			localConn = nil
			_ = conn.Close()
			localChan <- err
			continue
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
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

		dup := make([]byte, n)
		copy(dup, buf[:n])
		remoteQueue <- dup
	}
}

func handleRemote(localAddr, remoteAddr string, macChan, macChanRemote chan string, remoteConn net.Conn, remoteConnChan chan net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	defer close(remoteChan)
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v][%v]remote  ", mac, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	go func() {
		for m := range macChanRemote {
			if m != mac {
				logger.SetPrefix(fmt.Sprintf("[%17v][%v]remote  ", m, remoteAddr))
			}
			mac = m
		}
	}()

	go func() {
		for bb := range remoteQueue {
			conn := remoteConn
			if nil == conn {
				continue
			}
			_, err = conn.Write(bb)
			if err != nil {
			} // 发送数据
			if len(bb) > 8 && 0xFE == bb[0] && 0x01 == bb[3] {
				if 0x31 == bb[7] || 0x34 == bb[7] || 0x35 == bb[7] {
					logger.Printf("Write %#x with %v\n", bb[7], err)
				}
			}
		}
	}()

	reader := bufio.NewReader(remoteConn)
	var buf [1024]byte
	var n int
	for {
		if nil == remoteConn {
			for {
				if nil == remoteConn {
					remoteConn = <-remoteConnChan
				}
				if nil == remoteConn {
					remoteChan <- errors.New("need retry")
					continue
				}
				// reader.Reset(remoteConn)
				reader = bufio.NewReader(remoteConn)
				break
			}
		}
		n, err = reader.Read(buf[:]) // todo 粘包
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Printf("Connection Read EOF %v\n", err)
			} else {
				logger.Printf("Connection Read failed %v\n", err)
			}
			conn := remoteConn
			remoteConn = nil
			_ = conn.Close()
			remoteChan <- err
			continue
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x31 == buf[7] || 0x34 == buf[7] || 0x35 == buf[7] {
				logger.Printf("Read %#v\n", buf[7])
			}
		}

		dup := make([]byte, n)
		copy(dup, buf[:n])
		localQueue <- dup
	}
}
