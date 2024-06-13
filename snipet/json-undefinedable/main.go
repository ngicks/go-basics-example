package main

import (
	"encoding/json"
	"fmt"
)

type opt[V any] struct {
	ok bool
	v  V
}

func (o opt[V]) MarshalJSON() ([]byte, error) {
	if !o.ok {
		return []byte("null"), nil
	}
	return json.Marshal(o.v)
}

func (o *opt[V]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && string(data) == "null" {
		var zero V
		o.ok = false
		o.v = zero
		return nil
	}
	err := json.Unmarshal(data, &o.v)
	o.ok = err == nil
	return err
}

type und[T any] []opt[T]

func Undefined[T any]() und[T] {
	return nil
}

func Null[T any]() und[T] {
	return und[T]{opt[T]{ok: false}}
}

func Defined[T any](t T) und[T] {
	return und[T]{opt[T]{ok: true, v: t}}
}

func (u und[T]) MarshalJSON() ([]byte, error) {
	if len(u) == 0 {
		return []byte("null"), nil
	}
	return json.Marshal(u[0])
}

func (u *und[T]) clean() {
	uu := (*u)[:cap(*u)]
	for i := 0; i < len(uu); i++ {
		// shouldn't be happening,
		// at least erase them
		// so that held items are
		// up to GC.
		uu[i] = opt[T]{}
	}
}

func (u *und[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		if len(*u) == 0 {
			*u = []opt[T]{{}}
			return nil
		}
		u.clean()
		*u = (*u)[:1]
		(*u)[0] = opt[T]{}
		return nil
	}
	var o opt[T]
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}
	if len(*u) == 0 {
		*u = append(*u, o)
	} else {
		u.clean()
		(*u)[0] = o
	}
	return nil
}

type Und struct {
	Foo string      `json:",omitempty"`
	Bar und[string] `json:",omitempty"`
}

func main() {
	for _, v := range []Und{
		{"foo", Undefined[string]()},
		{"foo", Null[string]()},
		{"foo", Defined("bar")},
	} {
		bin, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		fmt.Printf("marshaled = %s\n", bin)
		var u Und
		err = json.Unmarshal(bin, &u)
		if err != nil {
			panic(err)
		}
		fmt.Printf("unmarshaled = %#v\n", u)
	}
}
