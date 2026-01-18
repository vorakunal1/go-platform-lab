package main

import "fmt"

func increment(num *int, s string) {
	fmt.Println("inside ", s);
	for i:=0; i<10000; i++ {
	fmt.Println("inside ", s, i);
		*num++;
	}
}

func main() {
	num := 0;
	go increment(&num, "a")
	go increment(&num, "c")
	go increment(&num, "q")
	go increment(&num, "e")
	go increment(&num, "r")
	go increment(&num, "t")
	go increment(&num, "y")
	go increment(&num, "u")
	go increment(&num, "pu")
	increment(&num, "b")
	fmt.Println(num)
}