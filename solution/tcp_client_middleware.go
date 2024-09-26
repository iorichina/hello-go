package main

import (
	"bufio"
	"errors"
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
	logger := log.New(os.Stdout, "main   ", log.Ldate|log.Lmicroseconds)

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
		log.Fatalf("local ip-port %v unreachable err %v", localAddr, err.Error())
	}
	logger.Printf("[%v]local Connection Connected\n", localAddr)

	remoteTimeoutDuration := time.Duration(remoteTimeout) * time.Millisecond
	remoteConn, err := net.DialTimeout("tcp", remoteAddr, remoteTimeoutDuration)
	if nil != err {
		log.Fatalf("remote ip-port %v unreachable err %v", remoteAddr, err.Error())
	}
	logger.Printf("[%v]remote Connection Connected\n", remoteAddr)

	clientMiddleware(localAddr, localTimeoutDuration, remoteAddr, remoteTimeoutDuration, localConn, remoteConn)
}

func clientMiddleware(localAddr string, localTimeoutDuration time.Duration, remoteAddr string, remoteTimeoutDuration time.Duration, localConn, remoteConn net.Conn) {
	logger := log.New(os.Stdout, "middle ", log.Ldate|log.Lmicroseconds)
	var mac string
	var err error

	defer func(conn net.Conn) {
		err := conn.Close()
		logger.Printf("[%v]->[%v][%v]Connection local Close %v\n", localAddr, remoteAddr, mac, err)
	}(localConn) // 关闭连接
	defer func(conn net.Conn) {
		err := conn.Close()
		logger.Printf("[%v]->[%v][%v]Connection remote Close %v\n", localAddr, remoteAddr, mac, err)
	}(remoteConn) // 关闭连接

	localChan := make(chan error)
	defer close(localChan)
	remoteChan := make(chan error)
	defer close(remoteChan)

	localQueue := make(chan []byte, 128)
	defer close(localQueue)
	remoteQueue := make(chan []byte, 128)
	defer close(remoteQueue)

	macChan := make(chan string, 128)
	defer close(macChan)
	go func() {
		for m := range macChan {
			mac = m
		}
	}()

	localConnChan := make(chan net.Conn)
	defer close(localConnChan)
	remoteConnChan := make(chan net.Conn)
	defer close(remoteConnChan)

	go handleLocal(localAddr, remoteAddr, macChan, localConn, localConnChan, localQueue, remoteQueue, localChan, remoteChan)
	go handleRemote(localAddr, remoteAddr, macChan, remoteConn, remoteConnChan, localQueue, remoteQueue, localChan, remoteChan)

	for {
		select {
		case err = <-localChan:
			start := time.Now()
			logger.Printf("[%v][%v]Connection local err, %v, trying to reconnect\n", localAddr, mac, err)
			localConn, err = net.DialTimeout("tcp", localAddr, localTimeoutDuration)
			logger.Printf("[%v][%v]Connection local reconnect %v with %v\n", localAddr, mac, localConn, err)
			duration := localTimeoutDuration - time.Since(start)
			if nil != err && duration > 0 {
				time.Sleep(duration)
			}
			localConnChan <- localConn
		case err = <-remoteChan:
			start := time.Now()
			logger.Printf("[%v][%v]Connection remote err, %v, trying to reconnect\n", remoteAddr, mac, err)
			remoteConn, err = net.DialTimeout("tcp", remoteAddr, remoteTimeoutDuration)
			logger.Printf("[%v][%v]Connection remote reconnect %v with %v\n", remoteAddr, mac, remoteConn, err)
			duration := remoteTimeoutDuration - time.Since(start)
			if nil != err && duration > 0 {
				time.Sleep(duration)
			}
			remoteConnChan <- remoteConn
		}
	}
}

func handleLocal(localAddr, remoteAddr string, macChan chan string, localConn net.Conn, localConnChan chan net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	logger := log.New(os.Stdout, "local  ", log.Ldate|log.Lmicroseconds)
	var mac string
	var err error

	reader := bufio.NewReader(localConn)
	var buf [1024]byte
	var n int

	go func() {
		for bb := range localQueue {
			conn := localConn
			if nil == conn {
				if len(bb) > 8 && 0xFE == bb[0] && 0x01 == bb[3] {
					if 0x31 == bb[7] {
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
				logger.Printf("[%v][%v]Write err %v\n", localAddr, mac, err)
			}
		}
	}()

	_, err = localConn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
	logger.Printf("[%v][%v]Write(0x34) with %v\n", localAddr, mac, err)
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
				logger.Printf("[%v][%v]Write(0x34) with %v\n", localAddr, mac, err)
				break
			}
		}
		n, err = reader.Read(buf[:]) // todo 粘包
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Printf("[%v][%v]Connection Read EOF %v\n", localAddr, mac, err)
			} else {
				logger.Printf("[%v][%v]Connection Read failed %v\n", localAddr, mac, err)
			}
			conn := localConn
			localConn = nil
			_ = conn.Close()
			localChan <- err
			continue
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x34 == buf[7] {
				mac = strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
				macChan <- mac
				logger.Printf("[%v][%v]status query resp %#x\n", localAddr, mac, buf[8])
			} else if 0x35 == buf[7] {
				mac = strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
				macChan <- mac
				logger.Printf("[%v][%v]heartbeat\n", localAddr, mac)
			} else {
				logger.Printf("[%v][%v]%#v\n", localAddr, mac, buf[7])
			}
		}

		dup := make([]byte, n)
		copy(dup, buf[:n])
		remoteQueue <- dup
	}
}

func handleRemote(localAddr, remoteAddr string, macChan chan string, remoteConn net.Conn, remoteConnChan chan net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	logger := log.New(os.Stdout, "remote ", log.Ldate|log.Lmicroseconds)
	var mac string
	var err error

	go func() {
		for m := range macChan {
			mac = m
		}
	}()

	reader := bufio.NewReader(remoteConn)
	var buf [1024]byte
	var n int

	go func() {
		for bb := range remoteQueue {
			conn := remoteConn
			if nil == conn {
				continue
			}
			_, err = conn.Write(bb)
			if err != nil {
				logger.Printf("[%v][%v]Write err %v\n", remoteAddr, mac, err.Error())
			} // 发送数据
		}
	}()

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
				logger.Printf("[%v][%v]Connection Read EOF %v\n", remoteAddr, mac, err)
			} else {
				logger.Printf("[%v][%v]Connection Read failed %v\n", remoteAddr, mac, err)
			}
			conn := remoteConn
			remoteConn = nil
			_ = conn.Close()
			remoteChan <- err
			continue
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x34 == buf[7] {
				logger.Printf("[%v][%v]status query\n", remoteAddr, mac)
			} else {
				logger.Printf("[%v][%v]%#v\n", remoteAddr, mac, buf[7])
			}
		}

		dup := make([]byte, n)
		copy(dup, buf[:n])
		localQueue <- dup
	}
}
