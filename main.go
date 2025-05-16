package main

import (
	"crypto/rand"
	"github.com/axiomhq/hyperloglog"
	"runtime/trace"
	"sync"
)

type void struct{}

func (v *void) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func main() {
	var v void
	if err := trace.Start(&v); err != nil {
		panic(err)
	}
	defer trace.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go hll(&wg)
	}
	wg.Wait()
}

//go:noinline
func hll(wg *sync.WaitGroup) {
	buf := make([]byte, 1000)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	h := hyperloglog.New14()
	for i := 0; i < 100_000_000; i++ {
		// Allocate some memory. The SIGSEGV does not seem to reproduce without
		// this.
		alloc := make([]byte, 1000)
		copy(alloc, buf)
		h.Insert(buf)
	}
	wg.Done()
}
