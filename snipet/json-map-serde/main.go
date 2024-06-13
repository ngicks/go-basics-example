package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Printf("using map[string]any:\n")
	for _, bin := range [][]byte{
		[]byte(`{"foo":"bar", "baz":[1,2,3]}`),
		[]byte(`{"foo":"bar", "baz":[1,2,3], "qux": {"nested":"nested", "null":null}}`),
	} {
		m := make(map[string]any)
		err := json.Unmarshal(bin, &m)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %#v\n", m)
		/*
			map[string]interface {}{"baz":[]interface {}{1, 2, 3}, "foo":"bar"}
			map[string]interface {}{"baz":[]interface {}{1, 2, 3}, "foo":"bar", "qux":map[string]interface {}{"nested":"nested", "null":interface {}(nil)}}
		*/
	}

	fmt.Printf("using any:\n")
	for _, litBin := range [][]byte{
		[]byte(`123`),
		[]byte(`0.4`),
		[]byte(`true`),
		[]byte(`null`),
		[]byte(`["yay", 123]`),
		[]byte(`{"object":"yes"}`),
		[]byte(`"nay"`),
	} {
		var m any
		err := json.Unmarshal(litBin, &m)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %#v\n", m)
		/*
			123
			0.4
			true
			<nil>
			[]interface {}{"yay", 123}
			map[string]interface {}{"object":"yes"}
			"nay"
		*/
	}
}
