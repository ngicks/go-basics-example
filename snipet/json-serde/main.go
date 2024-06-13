package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

type Sample struct {
	Foo string
	Bar int
}

func main() {
	os.OpenFile("/path/to/file", os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_EXCL, fs.ModePerm)
	// []byte, errorを返す。
	bin, err := json.Marshal(Sample{"foo", 123})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(bin)) // {"Foo":"foo","Bar":123}

	var s Sample
	err = json.Unmarshal(bin, &s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", s) // {Foo:foo Bar:123}
}
