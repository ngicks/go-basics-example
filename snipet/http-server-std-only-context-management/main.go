package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := http.NewServeMux()

	type keyTy string
	var (
		baseKey keyTy = "base-key"
		connKey keyTy = "conn-key"
	)
	server := &http.Server{
		Handler: mux,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, connKey, "conn")
		},
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(context.Background(), baseKey, "base")
		},
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"context",
			slog.Any(string(baseKey), r.Context().Value(baseKey)),
			slog.Any(string(connKey), r.Context().Value(connKey)),
		)
	}))

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("listening = %s", listener.Addr()))
	logger.Info(fmt.Sprintf("server closed = %v", server.Serve(listener)))
}
