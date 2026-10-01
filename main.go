package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

const maxLineSize = 1024 * 1024

func handleConn(conn net.Conn, store *SafeMap) {
	addr := conn.RemoteAddr().String()

	log.Printf("connect: %s", addr)
	defer log.Printf("disconnect: %s", addr)
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 4096), maxLineSize)
	for scanner.Scan() {
		req, err := ParseRequest(scanner.Text())
		if err != nil {
			WriteErr(conn, err.Error())
			continue
		}

		result, err := commands[req.Command].Handler(store, req.Args)
		if err != nil {
			WriteErr(conn, err.Error())
		} else {
			WriteOK(conn, result)
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		WriteErr(conn, scanErr.Error())
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening on :8080")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("\nreceived interrupt, closing listener......")
		listener.Close()
	}()

	store := NewSafeMap()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go handleConn(conn, store)
	}
}
