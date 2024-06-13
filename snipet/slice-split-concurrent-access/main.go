package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	imageBuffer := make([]byte, 1<<30)

	bufChan := make(chan []byte)

	var wg sync.WaitGroup
	for range runtime.GOMAXPROCS(0) { // resource limit. in this case, max parallel computation (=CPU num)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case buf := <-bufChan:
					for i := range buf {
						buf[i] = byte(rand.N(255)) + 1
					}
				}
			}
		}()
	}

	for i := 0; i < len(imageBuffer); i += 64 * 1024 {
		upper := i + 64*1024
		if upper > len(imageBuffer) {
			upper = len(imageBuffer)
		}
		bufChan <- imageBuffer[i:upper]
	}

	cancel()
	wg.Wait()

	for i := range imageBuffer {
		if imageBuffer[i] == 0 {
			fmt.Println("invalid buffer content")
		}
	}
}
