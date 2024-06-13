package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func startTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bin, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for k, v := range r.Header {
			fmt.Printf("req header: %s = %s\n", k, strings.Join(v, ", "))
		}
		fmt.Printf("bytes read: %d\n", len(bin))
		fmt.Printf("body=\n%s\n", bin)

		if reqId := r.Header.Get("X-Request-Id"); reqId != "" {
			w.Header().Set("X-Request-Id", reqId)
		}
		w.WriteHeader(http.StatusOK)
		return
	}))
}

type keyTy string

const (
	RequestIdKey keyTy = "request-id"
)

type RequestIdTrapper struct {
	RoundTripper http.RoundTripper
}

func (t *RequestIdTrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	req, err := addRequestID(req)
	if err != nil {
		return nil, err
	}
	return t.RoundTripper.RoundTrip(req)
}

func addRequestID(req *http.Request) (*http.Request, error) {
	if reqId := req.Header.Get("X-Request-Id"); reqId != "" {
		return req, nil
	}

	// Clone req to avoid data race.
	req = req.Clone(req.Context())

	reqId, _ := req.Context().Value(RequestIdKey).(string)

	if reqId == "" {
		var buf [16]byte
		_, err := io.ReadFull(rand.Reader, buf[:])
		if err != nil {
			return nil, fmt.Errorf("generating random X-Request-Id: %w", err)
		}
		reqId = fmt.Sprintf("%x", buf)
	}

	req.Header.Set("X-Request-Id", reqId)

	return req, nil
}

func main() {
	server := startTestServer()
	defer server.Close()
	client := &http.Client{
		Transport: &RequestIdTrapper{RoundTripper: http.DefaultTransport},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		fmt.Printf("resp header: %s = %s\n", k, strings.Join(v, ", "))
	}

	fmt.Printf("\n\n")

	ctx = context.WithValue(ctx, RequestIdKey, "foobarbaz")
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		fmt.Printf("resp header: %s = %s\n", k, strings.Join(v, ", "))
	}
}
