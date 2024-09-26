package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/now"
	"hello-go/char"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Rectangle struct {
	width, height float64
}

type Circle struct {
	Rectangle
	radius float64
}

func (r Rectangle) area() float64 {
	return r.width * r.height
}

func (c Circle) area() float64 {
	return c.width * c.radius * math.Pi
}

func AfterFunc(f func()) {
	goFunc(f)
}
func goFunc(arg any) { //type any = interface{}
	switch arg.(type) {
	case int:
		fmt.Println(arg, "is an int value.")
	case string:
		fmt.Println(arg, "is a string value.")
	case int64:
		fmt.Println(arg, "is an int64 value.")
	case func():
		fmt.Println(arg, "is an func() type.")
	default:
		fmt.Println(arg, "is an unknown type.")
	}
	//强制转换 arg 为 func() 类型（无参数），再执行该 func
	go arg.(func())()
}

type Human interface {
	say(s string)
}

type Parent struct {
	Human
	s string
}

//func (p Parent) say(s string) {
//	fmt.Println(s)
//	p.s = s
//}

func (p *Parent) say(s string) {
	fmt.Println(s)
	p.s = s
}
func main() {
	func(human Human) {
		human.say("hello parent")
		fmt.Printf("human %#v\n", human)
	}(&Parent{})

	AfterFunc(func() {
		fmt.Println("I'm AfterFunc")
	})
	resp, err := http.Get("https://baidu.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	//b, err := io.ReadAllAndClose(resp.Body)
	//fmt.Println(string(b))
	io.Copy(os.Stdout, resp.Body)

	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	//Ctrl+D in linux or Ctrl+Z in windows to stop Scan
	var cc int
	for input.Scan() && cc < 10 {
		counts[input.Text()]++
		cc++
	}
	// NOTE: ignoring potential errors from input.Err()
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}

	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	c1 := Circle{r1, 10}
	c2 := Circle{r2, 25}

	fmt.Println("Area of r1 is: ", r1.area())
	fmt.Println("Area of r2 is: ", r2.area())
	fmt.Println("Area of c1 is: ", c1.area())
	fmt.Println("Area of c2 is: ", c2.area())

	fmt.Println(time.Now())
	fmt.Println("hello world", now.BeginningOfDay())
	env := map[string]any{
		"foo": 100,
		"ss":  strconv.Itoa('a'),
		"bar": md5.New().Sum([]byte("200")),
	}
	fmt.Println("env：", env)
	fmt.Println("请输入：")
	var buffer [512]byte
	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		fmt.Println("输入错误：", err)
		return
	}
	if n < 10 {
		fmt.Println("输入字符少于", 10, char.CharLen("输入字符少于"))
		return
	}
	s := string(buffer[:])
	fmt.Println(s)
	fmt.Println("len()", len(buffer))
	fmt.Println("len(回)", len("回"))
	fmt.Println()
	fmt.Print("按回车退出...")
	os.Stdin.Read(buffer[:])
}
