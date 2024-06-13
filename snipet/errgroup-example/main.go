package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type task int

func tasksNaiveGroup(ctx context.Context, tasks []task, work func(ctx context.Context, t task) error) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	sem := make(chan struct{}, 15)

	wg.Add(len(tasks))
	for _, t := range tasks {
		sem <- struct{}{}
		go func() {
			defer func() {
				wg.Done()
				<-sem
			}()
			e := work(ctx, t)
			if e != nil {
				cancel(e)
			}
		}()
	}
	wg.Wait()
	return context.Cause(ctx)
}

func tasksErrgroup(ctx context.Context, tasks []task, work func(ctx context.Context, t task) error) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(15)
	for _, t := range tasks {
		g.Go(func() error {
			return work(ctx, t)
		})
	}
	return g.Wait()
}

func tasksErrgroupRepanic(ctx context.Context, tasks []task, work func(ctx context.Context, t task) error) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(15)
	var (
		panicOnce sync.Once
		panicked  any
	)
	for _, t := range tasks {
		g.Go(func() error {
			var err error
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				var set bool
				panicOnce.Do(func() {
					set = true
					panicked = rec
				})
				if set {
					err = fmt.Errorf("panicked: %v", rec)
				}
			}()
			err = work(ctx, t)
			return err
		})
	}
	err := g.Wait()
	if panicked != nil {
		panic(panicked)
	}
	return err
}

func main() {
	ctx := context.Background()
	var (
		tasks   []task
		blocker chan struct{}
		wg      sync.WaitGroup
	)

	for i := range 20 {
		tasks = append(tasks, task(i))
	}
	blocker = make(chan struct{})

	fmt.Printf("native group:\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = tasksNaiveGroup(ctx, tasks, func(ctx context.Context, t task) error {
			fmt.Printf("tid = %d, ", t)
			<-blocker
			return nil
		})
	}()

	time.Sleep(time.Second)
	fmt.Printf("\nslept one sec\n")
	close(blocker)

	wg.Wait()

	blocker = make(chan struct{})

	fmt.Printf("\n\nerrgroup:\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = tasksErrgroup(ctx, tasks, func(ctx context.Context, t task) error {
			fmt.Printf("tid = %d, ", t)
			<-blocker
			return nil
		})
	}()

	time.Sleep(time.Second)
	fmt.Printf("\nslept one sec\n")
	close(blocker)

	wg.Wait()

	fmt.Printf("\n\nerrgroup:\n")

	func() {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				rec := recover()
				fmt.Printf("panicked = %v\n", rec)
			}()
			_ = tasksErrgroupRepanic(ctx, tasks, func(ctx context.Context, t task) error {
				panic("foobar!")
			})
		}()
		wg.Wait()
	}()
	/*
	   native group:
	   tid = 14, tid = 7, tid = 8, tid = 9, tid = 12, tid = 13, tid = 6, tid = 0, tid = 11, tid = 1, tid = 3, tid = 4, tid = 5, tid = 2, tid = 10,
	   slept one sec
	   tid = 15, tid = 16, tid = 17, tid = 19, tid = 18,

	   errgroup:
	   tid = 14, tid = 7, tid = 8, tid = 9, tid = 10, tid = 13, tid = 2, tid = 12, tid = 0, tid = 1, tid = 4, tid = 6, tid = 3, tid = 5, tid = 11,
	   slept one sec
	   tid = 15, tid = 16, tid = 19, tid = 17, tid = 18,

	   errgroup:
	   panicked = foobar!
	*/
}
