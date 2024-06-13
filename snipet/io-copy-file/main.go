package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	switchCh := make(chan struct{})
	go func() {
		<-switchCh
		for ctx.Err() == nil {
			for range 1_000_000_000_000 {
				// This should cause go runtime to emit SIGURG
			}
		}
	}()
	switchCh <- struct{}{}

	tmpDir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	err = os.WriteFile(filepath.Join(tmpDir, "src"), []byte("fobarbaz\nbazbazbaz"), fs.ModePerm)
	if err != nil {
		panic(err)
	}

	src, err := os.Open(filepath.Join(tmpDir, "src"))
	if err != nil {
		panic(err)
	}
	defer func() { fmt.Printf("%#v\n", src.Close()) }()

	dst, err := os.Create(filepath.Join(tmpDir, "dst"))
	if err != nil {
		panic(err)
	}
	defer func() { fmt.Printf("%#v\n", dst.Close()) }()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}

	cancel()
}
