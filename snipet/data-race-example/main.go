package main

import (
	"fmt"
	"runtime"
	"slices"
	"sync"
)

func main() {
	var a []int

	var wg sync.WaitGroup
	for range runtime.GOMAXPROCS(0) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 100 {
				a = append(a, i)
			}
		}()
	}
	wg.Wait()

	group := map[int]int{}
	var keys []int
	for _, num := range a {
		if !slices.Contains(keys, num) {
			keys = append(keys, num)
		}
		group[num] = group[num] + 1
	}
	slices.Sort(keys)
	for _, key := range keys {
		num, count := key, group[key]
		if count != runtime.GOMAXPROCS(0) {
			fmt.Printf(
				"data race caused invalid state at num = %d, count = %d, expected = %d\n",
				num, count, runtime.GOMAXPROCS(0),
			)
		}
	}
	fmt.Printf("result = %#v\n", group)
	/*
	   data race caused invalid state at num = 0, count = 2, expected = 24
	   data race caused invalid state at num = 1, count = 2, expected = 24
	   data race caused invalid state at num = 2, count = 2, expected = 24
	   data race caused invalid state at num = 3, count = 2, expected = 24
	   data race caused invalid state at num = 4, count = 2, expected = 24
	   data race caused invalid state at num = 5, count = 2, expected = 24
	   data race caused invalid state at num = 6, count = 2, expected = 24
	   data race caused invalid state at num = 7, count = 2, expected = 24
	   data race caused invalid state at num = 8, count = 2, expected = 24
	   data race caused invalid state at num = 9, count = 2, expected = 24
	   data race caused invalid state at num = 10, count = 2, expected = 24
	   data race caused invalid state at num = 11, count = 2, expected = 24
	   data race caused invalid state at num = 12, count = 2, expected = 24
	   data race caused invalid state at num = 13, count = 2, expected = 24
	   data race caused invalid state at num = 14, count = 2, expected = 24
	   data race caused invalid state at num = 15, count = 2, expected = 24
	   data race caused invalid state at num = 16, count = 2, expected = 24
	   data race caused invalid state at num = 17, count = 2, expected = 24
	   data race caused invalid state at num = 18, count = 2, expected = 24
	   data race caused invalid state at num = 19, count = 2, expected = 24
	   data race caused invalid state at num = 20, count = 2, expected = 24
	   data race caused invalid state at num = 21, count = 2, expected = 24
	   data race caused invalid state at num = 22, count = 2, expected = 24
	   data race caused invalid state at num = 23, count = 2, expected = 24
	   data race caused invalid state at num = 24, count = 2, expected = 24
	   data race caused invalid state at num = 25, count = 2, expected = 24
	   data race caused invalid state at num = 26, count = 2, expected = 24
	   data race caused invalid state at num = 27, count = 2, expected = 24
	   data race caused invalid state at num = 28, count = 2, expected = 24
	   data race caused invalid state at num = 29, count = 2, expected = 24
	   data race caused invalid state at num = 30, count = 2, expected = 24
	   data race caused invalid state at num = 31, count = 2, expected = 24
	   data race caused invalid state at num = 32, count = 2, expected = 24
	   data race caused invalid state at num = 33, count = 2, expected = 24
	   data race caused invalid state at num = 34, count = 2, expected = 24
	   data race caused invalid state at num = 35, count = 2, expected = 24
	   data race caused invalid state at num = 36, count = 2, expected = 24
	   data race caused invalid state at num = 37, count = 2, expected = 24
	   data race caused invalid state at num = 38, count = 2, expected = 24
	   data race caused invalid state at num = 39, count = 2, expected = 24
	   data race caused invalid state at num = 40, count = 2, expected = 24
	   data race caused invalid state at num = 41, count = 2, expected = 24
	   data race caused invalid state at num = 42, count = 2, expected = 24
	   data race caused invalid state at num = 43, count = 2, expected = 24
	   data race caused invalid state at num = 44, count = 2, expected = 24
	   data race caused invalid state at num = 45, count = 2, expected = 24
	   data race caused invalid state at num = 46, count = 2, expected = 24
	   data race caused invalid state at num = 47, count = 2, expected = 24
	   data race caused invalid state at num = 48, count = 2, expected = 24
	   data race caused invalid state at num = 49, count = 2, expected = 24
	   data race caused invalid state at num = 50, count = 2, expected = 24
	   data race caused invalid state at num = 51, count = 2, expected = 24
	   data race caused invalid state at num = 52, count = 2, expected = 24
	   data race caused invalid state at num = 53, count = 2, expected = 24
	   data race caused invalid state at num = 54, count = 2, expected = 24
	   data race caused invalid state at num = 55, count = 2, expected = 24
	   data race caused invalid state at num = 56, count = 2, expected = 24
	   data race caused invalid state at num = 57, count = 2, expected = 24
	   data race caused invalid state at num = 58, count = 2, expected = 24
	   data race caused invalid state at num = 59, count = 2, expected = 24
	   data race caused invalid state at num = 60, count = 2, expected = 24
	   data race caused invalid state at num = 61, count = 2, expected = 24
	   data race caused invalid state at num = 62, count = 2, expected = 24
	   data race caused invalid state at num = 63, count = 2, expected = 24
	   data race caused invalid state at num = 64, count = 2, expected = 24
	   data race caused invalid state at num = 65, count = 2, expected = 24
	   data race caused invalid state at num = 66, count = 2, expected = 24
	   data race caused invalid state at num = 67, count = 2, expected = 24
	   data race caused invalid state at num = 68, count = 2, expected = 24
	   data race caused invalid state at num = 69, count = 2, expected = 24
	   data race caused invalid state at num = 70, count = 2, expected = 24
	   data race caused invalid state at num = 71, count = 2, expected = 24
	   data race caused invalid state at num = 72, count = 2, expected = 24
	   data race caused invalid state at num = 73, count = 2, expected = 24
	   data race caused invalid state at num = 74, count = 2, expected = 24
	   data race caused invalid state at num = 75, count = 2, expected = 24
	   data race caused invalid state at num = 76, count = 2, expected = 24
	   data race caused invalid state at num = 77, count = 2, expected = 24
	   data race caused invalid state at num = 78, count = 2, expected = 24
	   data race caused invalid state at num = 79, count = 2, expected = 24
	   data race caused invalid state at num = 80, count = 2, expected = 24
	   data race caused invalid state at num = 81, count = 2, expected = 24
	   data race caused invalid state at num = 82, count = 2, expected = 24
	   data race caused invalid state at num = 83, count = 2, expected = 24
	   data race caused invalid state at num = 84, count = 2, expected = 24
	   data race caused invalid state at num = 85, count = 2, expected = 24
	   data race caused invalid state at num = 86, count = 2, expected = 24
	   data race caused invalid state at num = 87, count = 2, expected = 24
	   data race caused invalid state at num = 88, count = 2, expected = 24
	   data race caused invalid state at num = 89, count = 2, expected = 24
	   data race caused invalid state at num = 90, count = 2, expected = 24
	   data race caused invalid state at num = 91, count = 2, expected = 24
	   data race caused invalid state at num = 92, count = 2, expected = 24
	   data race caused invalid state at num = 93, count = 2, expected = 24
	   data race caused invalid state at num = 94, count = 2, expected = 24
	   data race caused invalid state at num = 95, count = 2, expected = 24
	   data race caused invalid state at num = 96, count = 2, expected = 24
	   data race caused invalid state at num = 97, count = 2, expected = 24
	   data race caused invalid state at num = 98, count = 2, expected = 24
	   data race caused invalid state at num = 99, count = 2, expected = 24
	   result = map[int]int{0:2, 1:2, 2:2, 3:2, 4:2, 5:2, 6:2, 7:2, 8:2, 9:2, 10:2, 11:2, 12:2, 13:2, 14:2, 15:2, 16:2, 17:2, 18:2, 19:2, 20:2, 21:2, 22:2, 23:2, 24:2, 25:2, 26:2, 27:2, 28:2, 29:2, 30:2, 31:2, 32:2, 33:2, 34:2, 35:2, 36:2, 37:2, 38:2, 39:2, 40:2, 41:2, 42:2, 43:2, 44:2, 45:2, 46:2, 47:2, 48:2, 49:2, 50:2, 51:2, 52:2, 53:2, 54:2, 55:2, 56:2, 57:2, 58:2, 59:2, 60:2, 61:2, 62:2, 63:2, 64:2, 65:2, 66:2, 67:2, 68:2, 69:2, 70:2, 71:2, 72:2, 73:2, 74:2, 75:2, 76:2, 77:2, 78:2, 79:2, 80:2, 81:2, 82:2, 83:2, 84:2, 85:2, 86:2, 87:2, 88:2, 89:2, 90:2, 91:2, 92:2, 93:2, 94:2, 95:2, 96:2, 97:2, 98:2, 99:2}
	*/
}
