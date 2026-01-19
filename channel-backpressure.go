package main

import "fmt"
import "time"
import "sync"

func workers(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range ch {
		fmt.Println("slow consumer ", val)
		time.Sleep(2*time.Second)
	}
}

func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i:=0; i<100000; i++ {
		ch <- i
	}

	close(ch)
}

func main() {
	ch := make(chan int, 1)
	var wg sync.WaitGroup 

	wg.Add(2)
	go workers(ch, &wg)
	go producer(ch, &wg)

	wg.Wait()
	fmt.Println(runtime.NumGoroutine())
}