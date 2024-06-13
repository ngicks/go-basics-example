package main

import (
	"chan-one-of/oneof"
	"context"
	"fmt"
	"math/rand/v2"
)

func main() {
	for _, num := range []int{1, 2, 4, 5, 8, 10, 16, 19} {
		ctx, cancel := context.WithCancel(context.Background())

		fmt.Printf("num = %d:\n\n", num)

		targets := make([]chan int, num)
		for i := range targets {
			targets[i] = make(chan int)
		}

		sender := sendOnly(targets)
		receiver := recvOnly(targets)

		received := map[int]bool{}
		done := make(chan struct{})
		go func() {
			for {
				v, chosen, ok := oneof.Recv(receiver, ctx.Done())
				fmt.Printf("Recv: value = %d, chosen = %d, received = %t\n", v, chosen, ok)
				if !ok {
					break
				}
				received[chosen] = true
			}
			close(done)
		}()

		sent, completed := oneof.SendEach(sender, func() int { return rand.N(32) }, ctx.Done())
		fmt.Printf("completed = %t, sent = %#v\n", completed, sent)

		cancel()
		<-done

		if !contains0ToN(received, num) {
			panic("implementation error")
		}
		fmt.Println()
	}
}

func sendOnly[T ~[]C, C ~(chan E), E any](s T) []chan<- E {
	chans := make([]chan<- E, len(s))
	for i := range s {
		chans[i] = s[i]
	}
	return chans
}

func recvOnly[T ~[]C, C ~(chan E), E any](s T) []<-chan E {
	chans := make([]<-chan E, len(s))
	for i := range s {
		chans[i] = s[i]
	}
	return chans
}

func contains0ToN(s map[int]bool, max int) bool {
	for i := range max {
		if !s[i] {
			return false
		}
	}
	return true
}
