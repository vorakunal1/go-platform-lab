package main

import (
	"fmt"
	"sync"
	"time"
	"context"
)

func worker(wg *sync.WaitGroup, ctx context.Context, work <-chan int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context canceled")
			return
		default:
			fmt.Print("def")
		}
	
		select {
		case wk, ok := <-work:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Print(wk)
		}
	}
}


func producer(work chan<- int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for i:=1; i<1000000; i++ {
		select {
		case <-ctx.Done():
			close(work)
			return
		default:
		work <- i
		}
	}
}


func main() {
	var wg sync.WaitGroup
	work := make(chan int, 2)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(12*time.Second, func() {
		cancel()
	})

	wg.Add(3)
	go worker(&wg, ctx, work)
	go worker(&wg, ctx, work)
	go producer(work, &wg, ctx)

	wg.Wait()
}