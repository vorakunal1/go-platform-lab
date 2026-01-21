package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan int) {
	defer wg.Done()

	for {
		x, ok := <-jobs
		if !ok {
			fmt.Println("jobs channel closed")
			return
		}
		fmt.Print(x)
	}

}

func worker1(ctx context.Context, wg *sync.WaitGroup, jobs <-chan int) {
	defer wg.Done()

	for {
		x := <-jobs
		if err := ctx.Err(); err != nil {
			fmt.Printf("Stopping work: %v\n", err)
			return
		}

		fmt.Print(x)
	}

}

func producer(jobs chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		jobs <- i
	}

	close(jobs)
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	time.AfterFunc(13*time.Second, func() {
		fmt.Println("total leaks", runtime.NumGoroutine())
		cancel()
	})

	jobs := make(chan int, 23)

	wg.Add(4)
	go worker(ctx, &wg, jobs)
	go worker(ctx, &wg, jobs)
	go worker1(ctx, &wg, jobs)
	go producer(jobs, &wg)

	// fmt.Println("total leaks", runtime.NumGoroutine())

	wg.Wait()
	fmt.Println("total leaks", runtime.NumGoroutine())

}
