package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func main() {
	fmt.Printf(
		"🐤< ｺﾝﾆﾁﾊ！ ₍₍⁽⁽ %v ₎₎⁾⁾\n ",
		strings.Join(
			lo.Shuffle([]string{"🐔", "🐣", "🐧", "🐓"}),
			"₎₎⁾⁾ ₍₍⁽⁽",
		),
	)
}
