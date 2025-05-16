package main

import (
	"crypto/rand"
	"fmt"
	"github.com/axiomhq/hyperloglog"
	"math/big"
	"os"
	"runtime/trace"
	"sync"
)

const (
	bufLen = 10_000
	concurrency = 100
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
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
	for i := 0; i < 1_000_000; i++ {
		start, err := rand.Int(rand.Reader, big.NewInt(500))
		if err != nil {
			panic(err)
		}
		end, err := rand.Int(rand.Reader, big.NewInt(500))
		if err != nil {
			panic(err)
		}
		h.Insert(buf[start.Int64() : start.Int64()+end.Int64()])
	}
	fmt.Println(h.Estimate())
	wg.Done()
}
