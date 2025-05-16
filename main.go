package main

import (
	"crypto/rand"
	"runtime/trace"
	"sync"
	metro "github.com/dgryski/go-metro"

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

	var u uint64
	for i := 0; i < 100_000_000; i++ {
		// Allocate some memory. The SIGSEGV does not seem to reproduce without
		// this.
		alloc := make([]byte, 1000)
		copy(alloc, buf)
		u = hash(buf)
	}
	_ = u
	wg.Done()
}

//go:noinline
func hash(b []byte) uint64 {
	return metro.Hash64(b, 1337)
}

