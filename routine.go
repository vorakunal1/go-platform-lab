package main

import "fmt"
import "sync"

func increment(num *int, s string, mu *sync.Mutex, wg *sync.WaitGroup) {
	fmt.Println("inside ", s)
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		// fmt.Println("inside ", s, i);
		mu.Lock()
		*num++
		mu.Unlock()
	}
}

func main() {
	num := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(10)

	go increment(&num, "a", &mu, &wg)
	go increment(&num, "c", &mu, &wg)
	go increment(&num, "q", &mu, &wg)
	go increment(&num, "e", &mu, &wg)
	go increment(&num, "r", &mu, &wg)
	go increment(&num, "t", &mu, &wg)
	go increment(&num, "y", &mu, &wg)
	go increment(&num, "u", &mu, &wg)
	go increment(&num, "p", &mu, &wg)
	go increment(&num, "b", &mu, &wg)
	wg.Wait()
	fmt.Println(num)
}
