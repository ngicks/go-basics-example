package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

const bufSize = 8 * 1024

var bytesPool = &sync.Pool{
	New: func() any {
		b := make([]byte, bufSize)
		return &b
	},
}

func getBytes() *[]byte {
	return bytesPool.Get().(*[]byte)
}

func putBytes(b *[]byte) {
	if b == nil || len(*b) != bufSize || cap(*b) != bufSize {
		// reject grown / shrunk
		return
	}
	bytesPool.Put(b)
}

var bufPool = &sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func getBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func putBuf(b *bytes.Buffer) {
	if b.Cap() > 64*1024 {
		// See https://golang.org/issue/23199
		return
	}
	b.Reset()
	bufPool.Put(b)
}

// prevent WriteTo from being used
type onlyReader struct {
	r io.Reader
}

func (r onlyReader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

// prevent ReadFrom from being used.
type onlyWriter struct {
	w io.Writer
}

func (w onlyWriter) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func main() {
	bytesSlice := getBytes()
	defer putBytes(bytesSlice)

	buf := getBuf()
	defer putBuf(buf)

	_, err := io.CopyBuffer(
		onlyWriter{buf},
		onlyReader{bytes.NewReader([]byte(`foobarbaz`))},
		*bytesSlice,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("output: %s\n", buf.String())                                  // output: foobarbaz
	fmt.Printf("buf: len: %d, cap: %d\n", len(*bytesSlice), cap(*bytesSlice)) // buf: len: 8192, cap: 8192
	fmt.Printf("content: %s\n", *bytesSlice)                                  // content: foobarbaz
}
