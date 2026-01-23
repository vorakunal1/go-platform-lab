package main

import(
	"fmt"
	"sync"
	"context"
	"time"
)

func worker(wg *sync.WaitGroup, ctx context.Context, work <-chan int, name string) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Parent conext canceled", name)
			return
		case wk, ok := <- work:
			if !ok {
				fmt.Println("channel closed", name)
				return
			}
			fmt.Print(wk)
		}
	}

}

func producer(wg *sync.WaitGroup, ctx context.Context, work chan<- int, name string) {
	defer wg.Done()

	for i:=0; i<100000000; i++ {
		select {
		case <-ctx.Done():
			close(work)
			fmt.Println("Parent conext canceled", name)
			return
		default:
			work<-i
		}
	}
	close(work)
}

func main() {
	ctxParent, cancelParent := context.WithCancel(context.Background())
	ctxChild, cancelChild := context.WithTimeout(ctxParent, 8*time.Second)

	time.AfterFunc(15*time.Second, func() {
		cancelParent()
		cancelChild()
	})

	work := make(chan int, 4)
	var wg sync.WaitGroup;

	wg.Add(4)
	go worker(&wg, ctxParent, work, "parent")
	go worker(&wg, ctxParent, work, "parent")
	go worker(&wg, ctxChild, work, "child")
	go producer(&wg, ctxParent, work, "parent")

	wg.Wait()
}