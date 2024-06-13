package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

type config struct {
	GOPATH     string    `env:"GOPATH"`
	SERVER_URL *url.URL  `env:"SERVER_URL,notEmpty"`
	T1         time.Time `env:"T1"`
	List       []string  `env:"LIST" envSeparator:":"`
}

func main() {
	fmt.Printf("$GOPATH = %q\n", os.Getenv("GOPATH"))
	v, ok := os.LookupEnv("NONEXISTENT")
	fmt.Printf("$NONEXISTENT = %q, found = %t\n", v, ok)

	os.Setenv("SERVER_URL", "https://exmaple.com")
	os.Setenv("T1", "2022-03-06T12:23:54+09:00")
	os.Setenv("LIST", "foo:bar:baz")
	var c config
	err := env.Parse(&c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOPATH = %s\n", c.GOPATH)
	fmt.Printf("SERVER_URL = %s\n", c.SERVER_URL)
	fmt.Printf("T1 = %#v\n", c.T1)
	fmt.Printf("LIST = %#v\n", c.List)
}
