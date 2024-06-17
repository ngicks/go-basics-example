package main

import (
	"fmt"
	"net"
	"net/http"

	"http-server-oapi-codegen/api"
	"http-server-oapi-codegen/server"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/oapi-codegen/echo-middleware"
)

func main() {
	e := echo.New()

	spec, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	e.Use(echomiddleware.OapiRequestValidator(spec))

	api.RegisterHandlersWithBaseURL(e, api.NewStrictHandler(server.New(), nil), "")

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
