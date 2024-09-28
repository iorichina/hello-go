package main

import (
	"fmt"
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

	clientMiddleware1(localAddr, localTimeoutDuration, remoteAddr, remoteTimeoutDuration, localConn, remoteConn)
	logger.Printf("([%v]->[%v])process stoped\n", localAddr, remoteAddr)
}

func clientMiddleware1(localAddr string, localTimeoutDuration time.Duration, remoteAddr string, remoteTimeoutDuration time.Duration, localConn, remoteConn net.Conn) {
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
	defer close(localChan)
	remoteChan := make(chan error)
	defer close(remoteChan)

	macChan := make(chan string, 128)
	defer close(macChan)
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

	go handleLocal1(localAddr, remoteAddr, macChan, macChanLocal, localConn, remoteConn, localChan)
	go handleRemote1(localAddr, remoteAddr, macChan, macChanRemote, localConn, remoteConn, remoteChan)

	select {
	case err = <-localChan:
		logger.Printf("Connection local err, %v, stop process\n", err)
		return
	case err = <-remoteChan:
		logger.Printf("Connection remote err, %v, stop process\n", err)
		return
	}
}

func handleLocal1(localAddr, remoteAddr string, macChan, macChanLocal chan string, localConn, remoteConn net.Conn, localChan chan error) {
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

	_, err = localConn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
	logger.Printf("Write(0x34) by connected with %v\n", err)
	reader := localConn //bufio.NewReader(localConn)
	var buf [1024]byte
	var n int
	for {
		n, err = reader.Read(buf[:])
		if err != nil {
			localChan <- err
			break
		}
		if n <= 0 {
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
		remoteConn.Write(dup)
	}
}

func handleRemote1(localAddr, remoteAddr string, macChan, macChanRemote chan string, localConn, remoteConn net.Conn, remoteChan chan error) {
	var mac string
	var err error
	logger := log.New(os.Stdout, fmt.Sprintf("[%17v][%v]remote ", mac, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)

	go func() {
		for m := range macChanRemote {
			if m != mac {
				logger = log.New(os.Stdout, fmt.Sprintf("[%17v][%v]remote ", m, remoteAddr), log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
			}
			mac = m
		}
	}()

	reader := remoteConn //bufio.NewReader(remoteConn)
	var buf [1024]byte
	var n int
	for {
		n, err = reader.Read(buf[:])
		if err != nil {
			remoteChan <- err
			break
		}
		if n <= 0 {
			continue
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x31 == buf[7] || 0x34 == buf[7] || 0x35 == buf[7] {
				logger.Printf("Read %#v(%d)\n", buf[7], int(buf[1])*256+int(buf[2]))
			}
		}

		dup := make([]byte, n)
		copy(dup, buf[:n])
		localConn.Write(dup)
	}
}
