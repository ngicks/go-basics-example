package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ngicks/go-common/contextkey"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqId := c.Request().Header.Get("X-Request-Id")
			if reqId == "" {
				var bytes [16]byte
				_, err := io.ReadFull(rand.Reader, bytes[:])
				if err != nil {
					return err
				}
				reqId = fmt.Sprintf("%x", bytes)
			}
			c.SetRequest(
				c.Request().WithContext(
					contextkey.WithSlogLogger(
						c.Request().Context(),
						baseLogger.With(slog.String("request-id", reqId)),
					),
				),
			)
			return next(c)
		}
	})

	// fallback先にio.Discardに書き込むloggerを用意しておくと、context.Contextにロガーがない時ログを出さないという決断ができます。
	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e.GET("/", func(c echo.Context) error {
		logger := contextkey.ValueSlogLoggerFallback(c.Request().Context(), nopLogger)
		logger.Info("request")
		return nil
	})

	server := &http.Server{
		Handler: e,
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening = %s\n", listener.Addr())
	fmt.Printf("server closed = %v\n", server.Serve(listener))
}
