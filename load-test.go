package main

import (
	"fmt"
	"runtime"
	"context"
	"sync"
	"net"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())

	time.AfterFunc(10*time.Second, func() {
		fmt.Println("before cancel", runtime.NumGoroutine())
		cancel()
	})

	l, _ := net.Listen("tcp", ":2000")
	fmt.Println("before start", runtime.NumGoroutine())

	wg.Add(1)
	go acceptConnection(l, &wg, ctx)

	for i:=0; i<5000; i++ {
		wg.Add(1)
		go client(&wg, ctx)
	}


	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("myMethod took %s\n", elapsed)
	fmt.Println("after cancel", runtime.NumGoroutine())
}

func acceptConnection(l net.Listener, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	go func() {
		<-ctx.Done()
		l.Close()
	} ()

	for {
		conn, err := l.Accept()

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Error accepting conenction",err)
		}

		wg.Add(1)
		go handleConnection(conn, wg, ctx)
	}
}

func handleConnection(c net.Conn, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	defer c.Close()

	go func() {
		<-ctx.Done()
		c.SetReadDeadline(time.Now())
	}()

	for {
	    buf := make([]byte, 4096)
		c.SetReadDeadline(time.Now().Add(3*time.Second))
		n,err := c.Read(buf)

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Error reading message", err)
		}

		fmt.Println("Received message:", string(buf[:n]))
	}
}

func client(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:2000")
	if err != nil {
		fmt.Println("Error connection server", err)
		return
	}
	defer conn.Close()

	_, err1 := conn.Write([]byte("Hello from client!!"))

	if err1!=nil {
		fmt.Println("Error sending data", err)
	}

	<-ctx.Done()
}
