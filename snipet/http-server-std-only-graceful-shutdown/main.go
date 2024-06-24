package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("blocking on context")
		<-r.Context().Done()
		logger.Info("context canceled", slog.Any("err", r.Context().Err()), slog.Any("cause", context.Cause(r.Context())))
	}))

	server := &http.Server{
		Handler: mux,
	}

	server.RegisterOnShutdown(func() {
		logger.Info("on shutdown")
	})

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info(fmt.Sprintf("listening = %s", listener.Addr()))
		// Serveでブロックしておく
		logger.Info(fmt.Sprintf("server closed = %v", server.Serve(listener)))
	}()

	// 別のgoroutineでsignal待ちに入る
	wg.Add(1)
	go func() {
		defer wg.Done()

		// signal.Notifyでsignalを待ち受ける。
		// とりあえずSIGINT, SIGTERMで事足りる。環境によって決める。
		sigChan := make(chan os.Signal, 10)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// signalされるまで待つ。このサーバーはsignalされる以外に終了する手段がない
		sig := <-sigChan
		// osによらないメッセージを得られるのでprintしてもよい
		logger.Warn(fmt.Sprintf("signaled with %q", sig))

		// Shutdownは新しいrequestの受付を止めたうえで、処理中のrequestのresponseが帰るまで待つ。
		// 上記Handler中ではブロックしたままなのでこのctxは必ずtimeoutする。
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Info("requesting shutting down the server")
		err := server.Shutdown(ctx)
		logger.Error("server shutdown error", slog.Any("err", err))

		// 渡したcontext.Contextのキャンセルによるエラーなのかは下記のようにするとわかる。
		// ただし、内部的に作成されたcontext.Contextのtimeoutやcancellationは検知できない。
		// あくまで、このctxが返したエラーなのかどうかだけ。
		// errors.Is(err, context.DeadlineExceeded)あるいは
		// errors.Is(err ,context.Canceled)のほうが良い時もある。
		if err != nil && errors.Is(err, ctx.Err()) {
			logger.Error("closing server")
			err := server.Close()
			logger.Error("server close error", slog.Any("err", err))
		}
	}()

	wg.Wait()
}
