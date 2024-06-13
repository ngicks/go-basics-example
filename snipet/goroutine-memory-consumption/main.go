package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

var (
	num = flag.Uint("n", 1_000_000, "num of goroutines")
)

func main() {
	fmt.Printf("pid = %d\n", os.Getpid())
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	var wg sync.WaitGroup
	for range *num {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
		}()
	}
	wg.Wait()
	<-ctx.Done()
}
