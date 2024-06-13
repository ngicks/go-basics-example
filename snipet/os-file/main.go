package main

import (
	"io"
	"io/fs"
	"os"
)

func main() {
	f, err := os.OpenFile("/path/to/file", os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_EXCL, fs.ModePerm)
	if err != nil {
		panic(err)
	}
	var _ io.Reader = f
	var _ io.Writer = f
	var _ io.Closer = f
	var _ io.Seeker = f
	var _ io.ReaderAt = f
	var _ io.WriterAt = f
}
