package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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

type keyTy string

const (
	RequestIdKey keyTy = "request-id"
)

func modifyConfig[T any](base T, fn func(c T) T) T {
	return fn(base)
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RequestIDWithConfig(modifyConfig(
		middleware.DefaultRequestIDConfig,
		func(conf middleware.RequestIDConfig) middleware.RequestIDConfig {
			conf.RequestIDHandler = func(ctx echo.Context, s string) {
				ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), RequestIdKey, s)))
			}
			return conf
		},
	)))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("received path = %s", c.Request().URL.Path)
			err := next(c)
			if err != nil {
				c.Logger().Errorf("err: %#v", err)
			} else {
				c.Logger().Infof("no error")
			}
			return err
		}
	})

	var store sync.Map
	g := e.Group("/pp")
	g.POST("/:key", func(c echo.Context) error {
		c.Logger().Infof("request id associated with ctx: %s", c.Request().Context().Value(RequestIdKey))
		key := c.Param("key")

		if !strings.HasPrefix(c.Request().Header.Get("Content-Type"), "application/json") {
			return c.JSON(http.StatusBadRequest, postResult{Key: key, Err: "non json content type"})
		}

		var s Sample
		err := c.Bind(&s)
		if err != nil {
			return c.JSON(http.StatusBadRequest, postResult{Key: key, Err: "binding: " + err.Error()})
		}
		if err := s.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, postResult{Key: key, Err: err.Error()})
		}

		prev, loaded := store.Swap(key, s)
		return c.JSON(http.StatusOK, postResult{Key: key, Prev: prev, Swapped: loaded, Result: "ok"})
	})
	g.GET("/:key", func(c echo.Context) error {
		key := c.Param("key")
		val, loaded := store.Load(key)
		if !loaded {
			return c.JSON(http.StatusNotFound, getResult{Key: key, Err: "not found"})
		}
		return c.JSON(http.StatusOK, getResult{Key: key, Value: val})
	})

	// e.Start("127.0.0.1:8080") // You can use echo's Start to let it handle listener and server!

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
