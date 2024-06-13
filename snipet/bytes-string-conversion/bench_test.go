package main

import (
	"fmt"
	"testing"
	"unsafe"
)

var (
	data = []byte("はろーhogehogeです。")
	want = string(data)
)

func Benchmark_Sprintf(b *testing.B) {
	var str string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = fmt.Sprintf("%s", data)
	}
	b.StopTimer()
	if str != want {
		b.Fatalf("not equal: want = %s, got = %s", want, str)
	}
}

func Benchmark_String(b *testing.B) {
	var str string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = string(data)
	}
	b.StopTimer()
	if str != want {
		b.Fatalf("not equal: want = %s, got = %s", want, str)
	}
}

func Benchmark_UnsafeCast(b *testing.B) {
	var str string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = *(*string)(unsafe.Pointer(&data))
	}
	b.StopTimer()
	if str != want {
		b.Fatalf("not equal: want = %s, got = %s", want, str)
	}
}

func Benchmark_UnsafeString(b *testing.B) {
	var str string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = unsafe.String(&data[0], len(data))
	}
	b.StopTimer()
	if str != want {
		b.Fatalf("not equal: want = %s, got = %s", want, str)
	}
}
