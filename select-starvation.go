package main

import (
	"fmt"
	"context"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, ctx context.Context, work <-chan int) {
	defer wg.Done()

	for {
		select {
		case wk := <-work:
			fmt.Print(wk)
		case <-ctx.Done():
			fmt.Println("Context canceled")
			return
		}
	} 

}

func producer(wg *sync.WaitGroup, work chan<- int) {
	defer wg.Done()

	for i:=1; i<10000000; i++{
		work <- i
	}

	close(work)
}


func main() {
	var wg sync.WaitGroup;
	ctx, cancel := context.WithCancel(context.Background());

	time.AfterFunc(2*time.Second, func() {
		cancel()
	})
	work := make(chan int)

	wg.Add(3)
	go worker(&wg, ctx, work)
	go worker(&wg, ctx, work)
	go producer(&wg, work)

	wg.Wait()
}