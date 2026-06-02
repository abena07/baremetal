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

func handleConn(conn net.Conn) {
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
		default:
			WriteOK(conn, "")
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go handleConn(conn)
	}
}
