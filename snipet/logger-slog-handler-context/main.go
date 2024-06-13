package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"sync"
	"testing/slogtest"
	"time"
)

type keyConverter struct {
	key     any
	convert func(v any) any
}

type groupAttr struct {
	name string
	attr []slog.Attr
}

type ctxHandler struct {
	inner        slog.Handler
	ctxGroupName string
	keyMapping   map[any]string
	keyOrder     []keyConverter
	groups       []groupAttr
}

var _ slog.Handler = (*ctxHandler)(nil)

func newCtxHandler(inner slog.Handler, ctxGroupName string, keyMapping map[any]string, keyOrder []keyConverter) (slog.Handler, error) {
	if len(keyMapping) != len(keyOrder) {
		return nil, errors.New("ctxHandler: mismatching len(keyMapping) and len(keyOrder)")
	}
	for _, k := range keyOrder {
		_, ok := keyMapping[k.key]
		if !ok {
			return nil, errors.New("ctxHandler: keyOrder contains unknown key")
		}
	}
	return &ctxHandler{
		inner:        inner,
		ctxGroupName: ctxGroupName,
		keyMapping:   keyMapping,
		keyOrder:     keyOrder,
		groups:       []groupAttr{{name: "top"}},
	}, nil
}

func (h *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ctxHandler) Handle(ctx context.Context, record slog.Record) error {
	var ctxAttrs []slog.Attr
	for _, k := range h.keyOrder {
		v := ctx.Value(k.key)
		if k.convert != nil {
			v = k.convert(v)
		}
		ctxAttrs = append(ctxAttrs, slog.Attr{Key: h.keyMapping[k.key], Value: slog.AnyValue(v)})
	}

	topGrAttrs := h.groups[0].attr
	if h.ctxGroupName == "" {
		topGrAttrs = append(ctxAttrs, slices.Clone(topGrAttrs)...)
	}

	// reordering attached attrs.
	//
	// ctx attrs alway come first.
	// Groups come latter.
	// Attrs attached to record will be added to last group if any,
	// otherwise will be attached back to the record.

	var attachedAttrs []slog.Attr
	record.Attrs(func(a slog.Attr) bool {
		attachedAttrs = append(attachedAttrs, a)
		return true
	})
	if len(attachedAttrs) > 0 {
		record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	}
	if h.ctxGroupName != "" {
		record.AddAttrs(slog.Attr{Key: h.ctxGroupName, Value: slog.GroupValue(ctxAttrs...)})
	}
	record.AddAttrs(topGrAttrs...)
	if len(h.groups) == 1 {
		record.AddAttrs(attachedAttrs...)
	}

	groups := h.groups[1:]
	if len(groups) > 0 {
		g := groups[len(groups)-1]
		groupAttr := slog.Attr{Key: g.name, Value: slog.GroupValue(append(attachedAttrs, g.attr...)...)}
		if len(groups) > 1 {
			groups = groups[:len(groups)-1]
			for i := len(groups) - 1; i >= 0; i-- {
				g := groups[i]
				groupAttr = slog.Attr{Key: g.name, Value: slog.GroupValue(append(slices.Clone(g.attr), groupAttr)...)}
			}
		}
		record.AddAttrs(groupAttr)
	}
	return h.inner.Handle(ctx, record)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	groups := slices.Clone(h.groups)
	g := groups[len(groups)-1]
	g.attr = append(slices.Clone(g.attr), attrs...)
	groups[len(groups)-1] = g
	return &ctxHandler{
		inner:        h.inner,
		ctxGroupName: h.ctxGroupName,
		keyMapping:   h.keyMapping,
		keyOrder:     h.keyOrder,
		groups:       groups,
	}
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &ctxHandler{
		inner:        h.inner,
		ctxGroupName: h.ctxGroupName,
		keyMapping:   h.keyMapping,
		keyOrder:     h.keyOrder,
		groups:       append(slices.Clone(h.groups), groupAttr{name: name}),
	}
}

type keyTy string

const (
	RequestIdKey keyTy = "request-id"
	SyncMapKey   keyTy = "sync-map"
)

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	wrapHandler := func(h slog.Handler, ctxName string) (slog.Handler, error) {
		return newCtxHandler(
			h,
			ctxName,
			map[any]string{
				RequestIdKey:       "request-key",
				SyncMapKey:         "values",
				"request-received": "request-received",
			},
			[]keyConverter{
				{key: RequestIdKey},
				{key: SyncMapKey, convert: func(v any) any {
					m, ok := v.(*sync.Map)
					if !ok {
						return nil
					}
					values := map[string]any{}
					m.Range(func(key, value any) bool {
						values[key.(string)] = value
						return true
					})
					return values
				}},
				{key: "request-received", convert: func(_ any) any {
					return time.Now()
				}},
			},
		)
	}

	randomId := hex.EncodeToString(must(io.ReadAll(io.LimitReader(rand.Reader, 16))))
	store := &sync.Map{}
	store.Store("foo", "foo")
	store.Store("bar", 123)
	store.Store("baz", struct {
		Key   string
		Value string
	}{"baz", "bazbaz"})

	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestIdKey, randomId)
	ctx = context.WithValue(ctx, SyncMapKey, store)

	for _, ctxGroupName := range []string{"", "ctx"} {
		logger := slog.New(
			must(
				wrapHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
					ctxGroupName,
				),
			),
		)
		logger.DebugContext(ctx, "yay", slog.String("yay", "yayay"))
		logger.With("foo", "bar").WithGroup("nah").With("why", "why not").DebugContext(ctx, "nay")
	}
	/*
	   {"time":"2024-06-13T15:20:39.178198009Z","level":"DEBUG","msg":"yay","request-key":"d2a4f0e9bc80f229a2e91abc46f4752e","values":{"bar":123,"baz":{"Key":"baz","Value":"bazbaz"},"foo":"foo"},"request-received":"2024-06-13T15:20:39.17820893Z","yay":"yayay"}
	   {"time":"2024-06-13T15:20:39.17829812Z","level":"DEBUG","msg":"nay","request-key":"d2a4f0e9bc80f229a2e91abc46f4752e","values":{"bar":123,"baz":{"Key":"baz","Value":"bazbaz"},"foo":"foo"},"request-received":"2024-06-13T15:20:39.178298952Z","foo":"bar","nah":{"why":"why not"}}
	   {"time":"2024-06-13T15:20:39.178314551Z","level":"DEBUG","msg":"yay","ctx":{"request-key":"d2a4f0e9bc80f229a2e91abc46f4752e","values":{"bar":123,"baz":{"Key":"baz","Value":"bazbaz"},"foo":"foo"},"request-received":"2024-06-13T15:20:39.178315423Z"},"yay":"yayay"}
	   {"time":"2024-06-13T15:20:39.178328608Z","level":"DEBUG","msg":"nay","ctx":{"request-key":"d2a4f0e9bc80f229a2e91abc46f4752e","values":{"bar":123,"baz":{"Key":"baz","Value":"bazbaz"},"foo":"foo"},"request-received":"2024-06-13T15:20:39.17832927Z"},"foo":"bar","nah":{"why":"why not"}}
	*/
	var buf bytes.Buffer
	handler := must(wrapHandler(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}), ""))

	err := slogtest.TestHandler(
		handler,
		func() []map[string]any {
			// fmt.Println(buf.String())
			var ms []map[string]any
			for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
				if len(line) == 0 {
					continue
				}
				var m map[string]any
				if err := json.Unmarshal(line, &m); err != nil {
					panic(err)
				}
				ms = append(ms, m)
			}
			return ms
		},
	)
	fmt.Printf("slogtest.TestHandler = %v\n", err)
	// slogtest.TestHandler = <nil>
}
