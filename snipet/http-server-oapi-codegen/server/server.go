package server

import (
	"context"
	"fmt"
	"http-server-oapi-codegen/api"
	"math/rand/v2"
	"sync"

	"github.com/oapi-codegen/nullable"
)

var _ api.StrictServerInterface = (*serverInterface)(nil)

type serverInterface struct {
	m *sync.Map
}

func New() api.StrictServerInterface {
	return &serverInterface{
		m: new(sync.Map),
	}
}

// (GET /foo)
func (s *serverInterface) GetFoo(ctx context.Context, request api.GetFooRequestObject) (api.GetFooResponseObject, error) {
	foos := make(map[string]api.Foo)
	s.m.Range(func(key, value any) bool {
		foos[key.(string)] = value.(api.Foo)
		return true
	})
	if len(foos) == 0 {
		var fooErr api.FooError
		err := fooErr.FromError1(api.Error1{Foo: nullable.NewNullableWithValue("yay")})
		if err != nil {
			fmt.Printf("err = %v\n", err)
			return nil, err
		}
		return api.GetFoo404JSONResponse(fooErr), nil
	}
	return api.GetFoo200JSONResponse(foos), nil
}

// (POST /foo)
func (s *serverInterface) PostFoo(ctx context.Context, request api.PostFooRequestObject) (api.PostFooResponseObject, error) {
	rand := rand.N(10)

	switch rand { //機嫌が悪いとエラー
	case 0:
		return api.PostFoo400JSONResponse(api.FooError2{
			Foo: nullable.NewNullableWithValue("yay"),
		}), nil
	case 1:
		return api.PostFoo400JSONResponse(api.FooError2{
			Bar: nullable.NewNullableWithValue("yay"),
		}), nil
	case 2:
		return api.PostFoo400JSONResponse(api.FooError2{
			Baz: nullable.NewNullableWithValue("yay"),
		}), nil
	default:
		s.m.Store(request.Body.Name, *request.Body)
		return api.PostFoo200JSONResponse(*request.Body), nil
	}
}
