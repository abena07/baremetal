package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func handleConn(conn net.Conn, store *SafeMap) {
	addr := conn.RemoteAddr().String()

	log.Printf("connect: %s", addr)
	defer log.Printf("disconnect: %s", addr)
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		req, err := ParseRequest(scanner.Text())
		if err != nil {
			WriteErr(conn, err.Error())
			continue
		}

		switch req.Command {
		case "PING":
			WriteOK(conn, "PONG")
		case "SET":
			store.Set(req.Args[0], req.Args[1])
			WriteOK(conn, "")
		case "GET":
			val, ok := store.Get(req.Args[0])
			if !ok {
				WriteErr(conn, "key not found")
			} else {
				WriteOK(conn, val)
			}
		case "DEL":
			_, ok := store.Get(req.Args[0])
			if !ok {
				WriteErr(conn, "key not found")
			} else {
				store.Delete(req.Args[0])
				WriteOK(conn, "")
			}
		case "LIST":
			WriteOK(conn, strings.Join(store.List(), "|"))
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
