package main

import (
	"fmt"
	"context"
	"time"
	"sync"
)

import _ "net/http/pprof"
import "net/http"


func worker(ctx context.Context, jobs <- chan int, wg *sync.WaitGroup) {
	defer wg.Done();

	for {
		select {
		case <- ctx.Done():
			fmt.Println("context cancelled")
			return 
		case job, ok := <- jobs:
			if !ok {
				fmt.Println("jobs channel closed")
				return 
			}

			fmt.Println("current job", job)
		}
	} 
}

func producer(jobs chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i:=0; i<10; i++ {
		jobs <- i
	}

	close(jobs)
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	time.AfterFunc(533*time.Second, func() {
		cancel()
	})

	jobs := make(chan int, 3)

	wg.Add(3)
	go worker(ctx, jobs, &wg)
	go worker(ctx, jobs, &wg)
	go worker(ctx, jobs, &wg)
	// go producer(jobs, &wg)

	http.ListenAndServe("localhost:6060", nil)
	time.Sleep(123*time.Second)

	wg.Wait()
	// fmt.Println(runtime.NumGoroutine())

}