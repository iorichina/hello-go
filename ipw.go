package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	for {
		run()
		time.Sleep(300 * time.Second) // 暂停 5 秒
	}
}
func run() {
	err, s := curl()
	if err != nil {
		return
	}
	writeTtring(s)
}
func curl() (error, string) {
	resp, err := http.Get("http://4.ipw.cn/")
	if err != nil {
		fmt.Println(err)
		return err, ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err, ""
	}
	return nil, string(body)
}
func writeTtring(s string) {
	file, e := os.OpenFile("ipw.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer file.Close()
	file.WriteString(s)
	fmt.Println(s)
}
