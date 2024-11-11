package io

import (
	"bufio"
	"compress/gzip"
	"errors"
	"hello-go/logger"
	"io"
	"os"
	"strings"
)

func ReadAllAndClose(r io.ReadCloser) ([]byte, error) {
	all, err := io.ReadAll(r)
	r.Close()
	return all, err
}

// 缓存式读取文件，降低gc，兼容gz文件
//
// @return *os.File 记得 defer close
func PlainTextReader(filename string) (*os.File, *bufio.Reader, error) {
	//打开文件
	open, err := os.Open(filename)
	if err != nil {
		logger.Errorf("文件打开报错 %v error %v", filename, err)
		return nil, nil, errors.New("文件打开报错 " + err.Error())
	}

	rd := bufio.NewReader(open)

	//单纯的txt文件
	if !strings.Contains(filename, ".gz") {
		return open, rd, nil
	}
	//兼容gzip文件
	gd, err := gzip.NewReader(rd)
	if err != nil {
		logger.Errorf("文件gzip打开报错 %v error %v", filename, err)
		return open, rd, errors.New("文件gzip打开报错 " + err.Error())
	}
	return open, bufio.NewReader(gd), nil
}
