package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sync"
)

var cache sync.Map

func loadImage(name string) (image.Image, error) {
	v, ok := cache.Load(name)
	if !ok {
		v, _ = cache.LoadOrStore(
			name,
			sync.OnceValues(func() (image.Image, error) {
				f, err := os.Open(filepath.Join("image", name))
				if err != nil {
					return nil, err
				}
				return png.Decode(f)
			}),
		)
	}
	return v.(func() (image.Image, error))()
}
