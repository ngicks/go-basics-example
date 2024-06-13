package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	fmt.Println("Hello world", Foo)
	fmt.Println(lo.Without([]string{"foo", "bar", "baz"}, "bar"))
}
