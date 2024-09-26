package io

import (
	"io"
)

func ReadAllAndClose(r io.ReadCloser) ([]byte, error) {
	all, err := io.ReadAll(r)
	r.Close()
	return all, err
}
