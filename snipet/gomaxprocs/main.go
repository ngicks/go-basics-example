package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	blockByTightLoop = flag.Bool("b", false, "")
	blockMain        = flag.Bool("m", false, "")
)

func main() {
	flag.Parse()
	fmt.Printf("block = %t\n", *blockByTightLoop)
	fmt.Printf("gomaxprocs is %d\n", runtime.GOMAXPROCS(0))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// flood blocking goroutine
	for range runtime.GOMAXPROCS(0) {
		go func() {
			if *blockByTightLoop {
				for range 10_000_000_000 {
					// block this goroutine by tight loop
				}
			} else {
				<-ctx.Done()
			}
		}()
	}

	go func() {
		c := make(chan os.Signal, 10)
		signal.Notify(c)
		limit := 5
		for {
			select {
			case <-ctx.Done():
				signal.Stop(c)
				return
			case sig := <-c:
				if limit > 0 {
					limit--
					/*
						signal received: "urgent I/O condition"
						signal received: "urgent I/O condition"
						signal received: "urgent I/O condition"
						signal received: "urgent I/O condition"
						signal received: "urgent I/O condition"
					*/
					fmt.Printf("signal received: %q\n", sig)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("tick tok")
			case <-ctx.Done():
				return
			}
		}
	}()

	if *blockMain {
		for range 10_000_000_000 {
			// block main goroutine by tight loop
		}
		fmt.Printf("exited busy loop\n")
	}

	<-ctx.Done()

	/*
	   gomaxprocs is 24
	   tick tok
	   tick tok
	   tick tok
	   tick tok
	   exited busy loop
	   tick tok
	*/
}
