package main

import (
	"fmt"
	"sync"
	"time"
	"context"
)


func performTask(ctx context.Context, id int) error {
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Printf(" [Helper] %d finished", id)
		return nil
	case <- ctx.Done():
		return ctx.Err()
	}
}

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Worker %d] Started\n", id)
	
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Worker %d] Stopping: %v\n", id, ctx.Err())
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			if err := performTask(ctx, job); err != nil {
				fmt.Printf("[Worker %d] Task %d error: %v\n", id, job, err)
			}
		}
	}
}

func producer(ctx context.Context, jobs chan<- int) {
	defer close(jobs)

	for i:=1; i<10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("[Producer] stopping job feed")
			return
		case jobs<-i:
			fmt.Printf("[Producer] sent job %d\n", i)
		}
	}
}

func main() {
	rootCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	jobs := make(chan int, 2)

	var wg sync.WaitGroup

	for w:=1; w<=3; w++ {
		wg.Add(1)
		go worker(rootCtx, w, jobs, &wg)
	}

	go producer(rootCtx, jobs)

	wg.Wait()
	fmt.Println("Main: All components shut down gracefully.")

}