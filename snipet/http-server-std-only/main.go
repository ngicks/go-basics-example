package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Sample struct {
	Foo string
	Bar int
}

func (s Sample) Validate() error {
	var builder strings.Builder
	if s.Foo == "" {
		builder.WriteString("missing Foo")
	}
	switch n := s.Bar; {
	case n == 0:
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("missing Bar")
	case n < 0:
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("negative Bar")
	case n > 250:
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("too large > 250 Bar")
	}
	if builder.Len() > 0 {
		return fmt.Errorf("validation error: %s", builder.String())
	}
	return nil
}

type getResult struct {
	Key   string
	Value any    `json:",omitempty"`
	Err   string `json:",omitempty"`
}

type postResult struct {
	Key     string
	Prev    any    `json:",omitempty"`
	Swapped bool   `json:",omitempty"`
	Result  string `json:",omitempty"`
	Err     string `json:",omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "" && r.URL.Path != "/" {
			fmt.Printf("unknown path = %s\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"err":"path not found"}` + "\n"))
			return
		}
		if r.Method != http.MethodHead && r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"err":"method not allowed"}` + "\n"))
			return
		}
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"foobarbaz"}` + "\n"))
	}))

	var store sync.Map
	mux.Handle("POST /pp/{key}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		key := r.PathValue("key")

		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			w.WriteHeader(http.StatusBadRequest)
			_ = enc.Encode(postResult{Key: key, Err: "non json content type"})
			return
		}
		dec := json.NewDecoder(io.LimitReader(r.Body, 64<<20)) // hard limit on 64MiB
		dec.DisallowUnknownFields()
		var s Sample
		if err := dec.Decode(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = enc.Encode(postResult{Key: key, Err: "bad request shape"})
			return
		}
		if dec.More() {
			w.WriteHeader(http.StatusBadRequest)
			_ = enc.Encode(postResult{Key: key, Err: "junk data after json value"})
			return
		}
		if err := s.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = enc.Encode(postResult{Key: key, Err: err.Error()})
			return
		}

		prev, loaded := store.Swap(key, s)
		w.WriteHeader(http.StatusOK)
		_ = enc.Encode(postResult{Key: key, Prev: prev, Swapped: loaded, Result: "ok"})
	}))
	mux.Handle("GET /pp/{key}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		key := r.PathValue("key")
		val, loaded := store.Load(key)
		if !loaded {
			w.WriteHeader(http.StatusNotFound)
			enc := json.NewEncoder(w)
			_ = enc.Encode(getResult{
				Key: key,
				Err: "not found",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(getResult{
			Key:   key,
			Value: val,
		})
	}))
	mux.Handle("/pp/{key}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"err":"method not allowed"}` + "\n"))
	}))

	server := &http.Server{
		Handler: mux,
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening = %s\n", listener.Addr())
	fmt.Printf("server closed = %v\n", server.Serve(listener))
}
