package main

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

var (
	ErrNotEligibleState = errors.New("not eligible state")
)

type state int

const (
	stateA state = iota + 1
	stateB
	stateC
)

type condWorker struct {
	s    state
	cond *sync.Cond
}

func newCondWorker() *condWorker {
	return &condWorker{
		s:    stateA,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (w *condWorker) changeState(s state) {
	w.cond.L.Lock()
	defer w.cond.L.Unlock()
	w.cond.Broadcast()
	w.s = s
}

func (w *condWorker) do(doIf state, waitIf func(state) bool, f func()) error {
	w.cond.L.Lock()
	defer w.cond.L.Unlock()
	if w.s == doIf {
		f()
		return nil
	}
	if waitIf == nil || !waitIf(w.s) {
		return ErrNotEligibleState
	}
	for {
		w.cond.Wait() // lock is freed while being blocked on Wait
		// lock is now held.
		if w.s == doIf {
			break
		}
		if !waitIf(w.s) {
			return ErrNotEligibleState
		}
	}
	f()
	return nil
}

type varSet struct {
	doIf   state
	waitIf []state
}

func main() {
	w := newCondWorker()
	sChan := make(chan state)
	doChan := make(chan varSet)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for s := range sChan {
			fmt.Printf("changing state to %d\n", s)
			w.changeState(s)
			fmt.Printf("changed state to %d\n", s)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range doChan {
			fmt.Printf("doing if %d, would wait if s is one of %v\n", v.doIf, v.waitIf)
			err := w.do(
				v.doIf,
				func(s state) bool { return slices.Contains(v.waitIf, s) },
				func() { fmt.Println("...working...") },
			)
			fmt.Printf("done with %v\n", err)
		}
	}()
	doChan <- varSet{doIf: stateA}
	/*
		doing if 1, would wait if s is one of []
		...working...
		done with <nil>
	*/
	doChan <- varSet{doIf: stateB}
	/*
		doing if 2, would wait if s is one of []
		done with not eligible state
	*/
	doChan <- varSet{doIf: stateB, waitIf: []state{stateA}}
	/*
		doing if 2, would wait if s is one of [1]
	*/
	fmt.Println("sleeping...")
	time.Sleep(time.Millisecond)
	fmt.Println("woke up!")
	sChan <- stateB
	/*
		changing state to 2
		changed state to 2
		...working...
		done with <nil>
	*/
	doChan <- varSet{doIf: stateC, waitIf: []state{stateB}}
	/*
		doing if 3, would wait if s is one of [2]
	*/
	fmt.Println("sleeping...")
	time.Sleep(time.Millisecond)
	fmt.Println("woke up!")
	sChan <- stateA
	/*
		changing state to 1
		changed state to 1
		done with not eligible state
	*/
	close(sChan)
	close(doChan)
	wg.Wait()
}
