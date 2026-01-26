package main

import (
	"fmt"
	"context"
	"net"
	"sync"
	"time"
	"runtime"
	// "bufio"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var wg sync.WaitGroup

	l, err := net.Listen("tcp", ":2000")

	if err!=nil {
		fmt.Println("Error creating server", err)
	}
	
	wg.Add(1)
	go acceptConnection(l, &wg, ctx)

	for i:=1; i<4; i++ {
		wg.Add(1)
		go client(&wg, ctx)
	}

	fmt.Println("before", runtime.NumGoroutine())

	wg.Wait()
	fmt.Println("after", runtime.NumGoroutine())
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
			fmt.Println("[Server] Accept error:", err)
		}

		wg.Add(1)
		go handleConnection(wg, conn, ctx)
	}
}

func handleConnection(wg *sync.WaitGroup, c net.Conn, ctx context.Context) {
	defer wg.Done()
	defer c.Close()
	fmt.Println("New Connection")

	c.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 4096)
	stop := context.AfterFunc(ctx, func() {
        c.SetReadDeadline(time.Now()) 
    })
	defer stop()

	for {
		err := c.SetReadDeadline(time.Now().Add(10 * time.Second))

		fmt.Println("read timer reset")

		m, err := c.Read(buf)

		if ctx.Err() != nil {
			return
		}

		if err!=nil {
			fmt.Println("Error reading data", err)
		}

		fmt.Println("Received message:", string(buf[:m]))
	}
}

func client(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:2000")
	if err != nil {
		fmt.Println("Error connection server", err)
	}
	defer conn.Close()

	_, err1 := conn.Write([]byte("Hello from client!!"))

	if err1!=nil {
		fmt.Println("Error sending data", err)
	}

	<-ctx.Done()
}