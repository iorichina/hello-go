package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "net"
    "strings"
    "time"
)

func tcp_client() {
    
}

func process(conn net.Conn) {
    var mac string
    defer func(conn net.Conn) {
        err := conn.Close()
        if err != nil {
            now := time.Now().Format("2006-01-02 15:04:05.000")
            fmt.Printf("%v[%v][%v] Close err %v\t%#v\n", now, conn.RemoteAddr(), mac, err.Error(), err)
        }
    }(conn) // 关闭连接
    remote := conn.RemoteAddr()
    now := time.Now().Format("2006-01-02 15:04:05.000")
    fmt.Printf("%v[%v] Connected to %v\n", now, remote, conn.LocalAddr())

    var err error
    _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
    if err != nil {
        fmt.Printf("%v[%v](0x34) Write err %v\n", now, remote, err)
    }

    var n int
    reader := bufio.NewReader(conn)
    var buf [512]byte
    for {
        now = time.Now().Format("2006-01-02 15:04:05.000")
        n, err = reader.Read(buf[:]) // 读取数据
        if err != nil {
            if errors.Is(err, io.EOF) {
                fmt.Printf("%v[%v][%v]Read EOF %#v\n", now, remote, mac, err)
            } else {
                fmt.Printf("%v[%v][%v]Read failed err %v\t%#v\n", now, remote, mac, err.Error(), err)
            }
            break
        }

        now = time.Now().Format("2006-01-02 15:04:05.000")
        if n > 8 && 0xFE == buf[0] && 0x01 == buf[3] {
            if 0x34 == buf[7] {
                mac = strings.Join([]string{string(buf[9:11]), string(buf[11:13]), string(buf[13:15]), string(buf[15:17]), string(buf[17:19]), string(buf[19:21])}, ":")
            }
            if 0x35 == buf[7] {
                mac = strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
                fmt.Printf("%v[%v][%v]读 %#x 心跳\n", now, remote, mac, buf[7])
            }
            if 0x30 == buf[7] || 0x3b == buf[7] {
                fmt.Printf("%v[%v][%v]读 %#x\t% X\n", now, remote, mac, buf[7], buf[:n])
            }
        }

        _, err = client.Write(buf[:n])
        if err != nil {
            fmt.Printf("%v[%v][%v] Write err %v\n", now, remote, mac, err)
            continue
        } // 发送数据
    }
}

func tcp_server() {
    listen, err := net.Listen("tcp", "0.0.0.0:80")
    if err != nil {
        fmt.Printf("%v Listen() failed, err %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)
        return
    }
    for {
        conn, err := listen.Accept() // 监听客户端的连接请求
        if err != nil {
            fmt.Printf("%v Accept() failed, err: %#v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)
            continue
        }
        go process(conn) // 启动一个goroutine来处理客户端的连接请求
    }
}

func main() {
    tcp_server()
}
