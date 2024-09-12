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

// go build tcp_server.go
//
// [nohup] sudo ./tcp_server >>log.log 2>&1 &
//
// TCP Server端测试
// 处理函数
func process(conn net.Conn) {
    defer func(conn net.Conn) {
        err := conn.Close()
        if err != nil {
            now := time.Now().Format("2006-01-02 15:04:05.000")
            fmt.Printf("%v[%v] Close err %#v\n", now, conn.RemoteAddr(), err)
        }
    }(conn) // 关闭连接
    remote := conn.RemoteAddr()
    now := time.Now().Format("2006-01-02 15:04:05.000")
    fmt.Printf("%v[%v] Connection to %v\n", now, remote, conn.LocalAddr())

    var err error
    _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 11, 48, 90, 0, 149})
    if err != nil {
        fmt.Printf("%v[%v](0x30) Write err %v\n", now, remote, err)
    }
    _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 51, 60})
    if err != nil {
        fmt.Printf("%v[%v](0x33) Write err %v\n", now, remote, err)
    }
    _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 66, 75})
    if err != nil {
        fmt.Printf("%v[%v](0x42) Write err %v\n", now, remote, err)
    }
    _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 52, 61})
    if err != nil {
        fmt.Printf("%v[%v](0x34) Write err %v\n", now, remote, err)
    }

    var mac string
    var n int
    for {
        now = time.Now().Format("2006-01-02 15:04:05.000")
        reader := bufio.NewReader(conn)
        var buf [512]byte
        n, err = reader.Read(buf[:]) // 读取数据
        if err != nil {
            if errors.Is(err, io.EOF) {
                fmt.Printf("%v[%v][%v]Read EOF %#v\n", now, remote, mac, err)
            } else {
                fmt.Printf("%v[%v][%v]Read failed err %#v\t%v\n", now, remote, mac, err, err.Error())
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
                fmt.Printf("%v[%v][%v]读 %#x 状态 %d (0=idle,1=playing,>1=error)\n", now, remote, mac, buf[7], buf[8])
                continue
            }
            if 0x35 == buf[7] {
                mac = strings.Join([]string{string(buf[8:10]), string(buf[10:12]), string(buf[12:14]), string(buf[14:16]), string(buf[16:18]), string(buf[18:20])}, ":")
                fmt.Printf("%v[%v][%v]读 %#x 心跳\n", now, remote, mac, buf[7])

                _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 9, 51, 60})
                if err != nil {
                    fmt.Printf("%v[%v][%v](0x33) Write err %v\n", now, remote, mac, err)
                }
                _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 11, 49, 0, 1, 61})
                if err != nil {
                    fmt.Printf("%v[%v][%v](0x31) Write err %v\n", now, remote, mac, err)
                    continue
                }
                _, err = conn.Write([]byte{254, 73, 66, 1, 182, 189, 11, 49, 0, 1, 61})
                if err != nil {
                    fmt.Printf("%v[%v][%v](0x30) Write err %v\n", now, remote, mac, err)
                    continue
                }
                _, err = conn.Write([]byte{254, 134, 226, 1, 121, 29, 10, 50, 2, 62})
                if err != nil {
                    fmt.Printf("%v[%v][%v](0x32) Write err %v\n", now, remote, mac, err)
                    continue
                }
                continue
            }
            fmt.Printf("%v[%v][%v]读 %#x\t% X\n", now, remote, mac, buf[7], buf[:n])
            continue
        }

        recvStr := string(buf[:n])
        fmt.Printf("%v[%v]读 %v\t% X\n", now, remote, recvStr, buf[:n])
        _, err = conn.Write([]byte(recvStr))
        if err != nil {
            fmt.Printf("%v[%v][%v] Write err %v\n", now, remote, mac, err)
            continue
        } // 发送数据
    }
}

func main() {
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
