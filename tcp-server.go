package main

import (
	"context"
	"fmt"
	"net"
	// "io"
	"bufio"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// import _ "net/http/pprof"
// import "net/http"

func acceptConnection(l net.Listener, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	fmt.Println("acceptCOnnection")

	go func() {
		<-ctx.Done()
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		// fmt.Println("inside acceptConnection", err)
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("[server] Accept error:", err)
		}

		wg.Add(1)
		go handleConnection(conn, wg, ctx)
	}
}

func handleConnection(c net.Conn, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	defer c.Close()

	reader := bufio.NewReader(c)
	fmt.Println("inside handleConnection", c)
	for {
		id, err := reader.ReadByte()
		if err != nil {
			return
		}
		fmt.Println("inside handleConnection1111111", id)
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		fmt.Printf("Received id: %d, Message: %s", id, message)

		if ctx.Err() != nil {
			return
		}
	}
}

func makeConnection(wg *sync.WaitGroup, counter *atomic.Int64, ctx context.Context) {
	defer wg.Done()
	defer counter.Add(-1)
	// var b byte = []byte(id)
	conn, err1 := net.Dial("tcp", "localhost:2000")
	defer conn.Close()

	if err1 != nil {
		fmt.Println("error connecting", err1)
		return
	}
	msg := []byte("Hello from the client!\n")
	// msg = append(b, msg)

	_, err := conn.Write(msg)
	if err != nil {
		fmt.Println("Error sending data:", err)
	}

	<-ctx.Done()
}

func main() {
	var wg sync.WaitGroup
	ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
	var connCount atomic.Int64
	l, err := net.Listen("tcp", ":2000")

	if err != nil {
		log.Fatal(err)
		fmt.Println("inside main", err)
	}

	defer l.Close()

	wg.Add(1)
	go acceptConnection(l, &wg, ctx)
	for i := 1; i < 4; i++ {
		wg.Add(1)
		connCount.Add(1)
		go makeConnection(&wg, &connCount, ctx)
	}

	fmt.Println("before",runtime.NumGoroutine())

	// http.ListenAndServe("localhost:6060", nil)

	wg.Wait()

	fmt.Println("after",runtime.NumGoroutine())
}
