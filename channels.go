package main

import "fmt"

func chan_util1(c chan int) {
	for i:=0; i<10; i++ {
		c <- i
	}
}

func chan_util2(c chan int) {
	for i:=10; i<20; i++ {
		c <- i
	}
}


func chan_util3(c chan int) {
	for i:=20; i<30; i++ {
		c <- i
	}
}

func chan_util4(c chan int) {
	for i:=0; i<30; i++ {
		c <- i
	}
}

func main() {
	unbuff_ch := make(chan int)
	buff_ch := make(chan int, 2)
	go chan_util1(unbuff_ch)
	go chan_util2(unbuff_ch)
	go chan_util3(unbuff_ch)
	go chan_util4(buff_ch)

	for j:=0; j<30; j++ {
		x := <- unbuff_ch
		fmt.Println("unbuffered", x)
	}

	for j:=0; j<30; j++ {
		x := <- buff_ch
		fmt.Println("buffered", x)
	}

}