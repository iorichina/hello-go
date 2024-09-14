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

func handleLocal(logger *log.Logger, macChan chan string, localConn net.Conn, localConnChan chan net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	localAddr := localConn.RemoteAddr()
	var mac string
	defer func() {
        if r := recover(); r != nil {
			now := time.Now().Format("2006-01-02 15:04:05.000")
            logger.Printf("%v[%v][%v]local Recovered from panic: %#v", now, localAddr, mac, r)
        }
    }()
	var err error
	var n int
	reader := bufio.NewReader(localConn)
	readerChan := make(chan *bufio.Reader)
	var buf [512]byte
	go func() {
		for bb := range localQueue {
			if nil == localConn {
				for {
					if nil == localConn {
						localConn = <-localConnChan
					}
					if nil == localConn {
						continue
					}
					readerChan <-bufio.NewReader(localConn)
				}
			}
			_, err = localConn.Write(bb)
			if err != nil {
				now = time.Now().Format("2006-01-02 15:04:05.000")
				logger.Printf("%v[%v][%v]local Write err %v\t%#v\n", now, localAddr, mac, err.Error(), err)
			} // 发送数据
		}
	}()
	
	_, err = localConn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
	if err != nil {
		now = time.Now().Format("2006-01-02 15:04:05.000")
		logger.Printf("%v[%v][%v]local Write(0x34) err %v\n", now, localAddr, mac, err)
	}

	for {
		if nil == reader {
			reader = <-readerChan
		}
		n, err = reader.Read(buf[:]) // 读取数据
		if err != nil {
			now = time.Now().Format("2006-01-02 15:04:05.000")
			if errors.Is(err, io.EOF) {
				logger.Printf("%v[%v][%v]local Read EOF %#v\n", now, localAddr, mac, err)
			} else {
				logger.Printf("%v[%v][%v]local Read failed err %v\t%#v\n", now, localAddr, mac, err.Error(), err)
			}
			_ = localConn.Close()
			localConn = nil
			reader = nil
			localChan <- err
			continue
		}

		if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
			if 0x34 == buf[7] {
				mac = strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
				macChan <- mac
			}
			if 0x35 == buf[7] {
				mac = strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
				macChan <- mac
			}
		}

		dup := make([]byte, n)
		copy(dup, buf[:n])
		remoteQueue <- dup
	}
}

func handleRemote(logger *log.Logger, macChan chan string, remoteConn net.Conn, remoteConnChan chan net.Conn, localQueue, remoteQueue chan []byte, localChan, remoteChan chan error) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	remoteAddr := remoteConn.RemoteAddr()
	var mac string
	defer func() {
        if r := recover(); r != nil {
			now := time.Now().Format("2006-01-02 15:04:05.000")
            logger.Printf("%v[%v][%v]remote Recovered from panic: %#v", now, remoteAddr, mac, r)
        }
    }()
	var err error
	var n int
	go func() {
		for m := range macChan {
			mac = m
		}
	}()
	reader := bufio.NewReader(remoteConn)
	readerChan := make(chan *bufio.Reader)
	var buf [512]byte
	go func() {
		for bb := range remoteQueue {
			if nil == remoteConn {
				for {
					if nil == remoteConn {
						remoteConn = <-remoteConnChan
					}
					if nil == remoteConn {
						continue
					}
					readerChan <-bufio.NewReader(remoteConn)
				}
			}
			_, err = remoteConn.Write(bb)
			if err != nil {
				now = time.Now().Format("2006-01-02 15:04:05.000")
				logger.Printf("%v[%v][%v]remote Write err %v\t%#v\n", now, remoteAddr, mac, err.Error(), err)
			} // 发送数据
		}
	}()
	for {
		if nil == reader {
			reader = <-readerChan
		}
		n, err = reader.Read(buf[:]) // 读取数据
		if err != nil {
			now = time.Now().Format("2006-01-02 15:04:05.000")
			if errors.Is(err, io.EOF) {
				logger.Printf("%v[%v][%v]remote Read EOF %#v\n", now, remoteAddr, mac, err)
			} else {
				logger.Printf("%v[%v][%v]remote Read failed err %v\t%#v\n", now, remoteAddr, mac, err.Error(), err)
			}
			_ = remoteConn.Close()
			remoteConn = nil
			reader = nil
			remoteChan <- err
			continue
		}

		dup := make([]byte, n)
		copy(dup, buf[:n])
		localQueue <- dup
	}
}

func clientMiddleware(logger *log.Logger, localAddr, remoteAddr string, localConn, remoteConn net.Conn) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	var mac string
	var err error

	defer func() {
        if r := recover(); r != nil {
			now := time.Now().Format("2006-01-02 15:04:05.000")
            logger.Printf("%v[%v]>[%v][%v]Recovered from panic: %#v", now, localAddr, remoteAddr, mac, r)
        }
    }()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			now := time.Now().Format("2006-01-02 15:04:05.000")
			logger.Printf("%v[%v]>[%v][%v]localConn Close err %v\t%#v\n", now, localAddr, remoteAddr, mac, err.Error(), err)
		}
	}(localConn) // 关闭连接
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			now := time.Now().Format("2006-01-02 15:04:05.000")
			logger.Printf("%v[%v]>[%v][%v]remoteConn Close err %v\t%#v\n", now, localAddr, remoteAddr, mac, err.Error(), err)
		}
	}(remoteConn) // 关闭连接

	localChan := make(chan error)
	defer close(localChan)
	remoteChan := make(chan error)
	defer close(remoteChan)

	localQueue := make(chan []byte)
	defer close(localQueue)
	remoteQueue := make(chan []byte)
	defer close(remoteQueue)

	macChan := make(chan string)
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

	go handleLocal(logger, macChan, localConn, localConnChan, localQueue, remoteQueue, localChan, remoteChan)
	go handleRemote(logger, macChan, remoteConn, remoteConnChan, localQueue, remoteQueue, localChan, remoteChan)

	for {
		select {
		case err = <-localChan:
			logger.Printf("%v[%v][%v]local err %v\t%#v disconnect, try to reconnect\n", now, localAddr,mac,  err.Error(), err)
			localConn, _ = net.Dial("tcp", localAddr)
			logger.Printf("%v[%v][%v]local Connection Connected\n", now, localAddr, mac)
			localConnChan <-localConn
		case err = <-remoteChan:
			logger.Printf("%v[%v][%v]remote err %v\t%#v disconnect, try to reconnect\n", now,  remoteAddr,mac, err.Error(), err)
			remoteConn, _ = net.Dial("tcp", remoteAddr)
			logger.Printf("%v[%v][%v]remote Connection Connected\n", now, remoteAddr, mac)
			remoteConnChan <-remoteConn
		}
	}
}

func main() {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	if len(os.Args) <= 2 {
		log.Fatalf("client_middleware.log need local ip-port and remote ip-port")
	}
	f, err := os.OpenFile("client_middleware.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatalf("create file client_middleware.log failed: %v", err)
	}
	logger := log.New(io.MultiWriter(os.Stdout, f), "", 0)

	localAddr := os.Args[1]
	remoteAddr := os.Args[2]

	localConn, err := net.Dial("tcp", localAddr)
	if nil != err {
		log.Fatalf("local ip-port %v unreachable err:%v\t%#v", localAddr, err.Error(), err)
	}
	now = time.Now().Format("2006-01-02 15:04:05.000")
	logger.Printf("%v[%v]local Connection Connected\n", now, localAddr)
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if nil != err {
		log.Fatalf("remote ip-port %v unreachable err:%v\t%#v", remoteAddr, err.Error(), err)
	}
	now = time.Now().Format("2006-01-02 15:04:05.000")
	logger.Printf("%v[%v]remote Connection Connected\n", now, remoteAddr)

	clientMiddleware(logger, localAddr, remoteAddr, localConn, remoteConn)
}
