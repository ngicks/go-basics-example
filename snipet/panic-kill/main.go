package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() {
		rec := recover()
		fmt.Printf("recovered = %v\n", rec)
	}()

	switcher := make(chan struct{})
	go func() {
		<-switcher
		panic("yay")
	}()
	switcher <- struct{}{}
	time.Sleep(time.Second)
	panic("nay")
}
