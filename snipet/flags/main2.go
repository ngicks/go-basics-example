package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	sub1 := flag.NewFlagSet("sub1", flag.PanicOnError)
	sub2 := flag.NewFlagSet("sub2", flag.PanicOnError)

	foo := sub1.String("foo", "", "foo")
	bar := sub2.Int("bar", 0, "bar")

	if len(os.Args) < 2 {
		panic("too short")
	}

	switch os.Args[1] {
	case "sub1":
		sub1.Parse(os.Args[2:])
		fmt.Printf("foo = %s\n", *foo)
	case "sub2":
		sub2.Parse(os.Args[2:])
		fmt.Printf("bar = %d\n", *bar)
	}
}
