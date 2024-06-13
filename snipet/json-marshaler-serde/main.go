package main

import (
	"encoding/json"
	"fmt"
)

type data1 struct {
	Foo string
}

type data2 struct {
	Bar int
}

type data3 struct {
	Baz bool
}

type Sample2 struct {
	tag  string
	data any
}

func (s Sample2) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	switch s.data.(type) {
	case data1:
		m["tag"] = "data1"
	case data2:
		m["tag"] = "data2"
	case data3:
		m["tag"] = "data3"
	default:
		return nil, fmt.Errorf("unknown error type")
	}
	bin, _ := json.Marshal(s.data)
	m["data"] = json.RawMessage(bin)

	return json.Marshal(m)
}

func (s *Sample2) UnmarshalJSON(data []byte) error {
	type T struct {
		Tag  string          `json:"tag"`
		Data json.RawMessage `json:"data"`
	}
	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	var raw any
	switch t.Tag {
	case "data1":
		var x data1
		err = json.Unmarshal(t.Data, &x)
		raw = x
	case "data2":
		var x data2
		err = json.Unmarshal(t.Data, &x)
		raw = x
	case "data3":
		var x data3
		err = json.Unmarshal(t.Data, &x)
		raw = x
	default:
		return fmt.Errorf("unknown tag")
	}
	if err != nil {
		return err
	}
	s.tag = t.Tag
	s.data = raw
	return nil
}

func main() {
	for _, d := range []Sample2{
		{data: data1{Foo: "foo"}},
		{data: data2{Bar: 5587}},
		{data: data3{Baz: true}},
	} {
		bin, err := json.Marshal(d)
		if err != nil {
			panic(err)
		}
		fmt.Printf("marshaled = %s\n", bin)
		var s Sample2
		err = json.Unmarshal(bin, &s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("unmarshaled = %#v\n", s)
		/*
			marshaled = {"data":{"Foo":"foo"},"tag":"data1"}
			unmarshaled = main.Sample2{tag:"data1", data:main.data1{Foo:"foo"}}
			marshaled = {"data":{"Bar":5587},"tag":"data2"}
			unmarshaled = main.Sample2{tag:"data2", data:main.data2{Bar:5587}}
			marshaled = {"data":{"Baz":true},"tag":"data3"}
			unmarshaled = main.Sample2{tag:"data3", data:main.data3{Baz:true}}
		*/
	}
}
