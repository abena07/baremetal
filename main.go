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

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line)

		req, errs := ParseRequest(line)

		if errs != nil {
			conn.Write([]byte(ResponseErrFormatter(errs.Error())))
		} else {
			conn.Write([]byte(ResponseOkFormatter(req.Command, req.Args)))
		}

		if line == "\r\n" {
			break
		}
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

	// req, err := ParseRequest("SET|1|adwoa\n")
	// if err != nil {
	// 	fmt.Println("error:", err)
	// } else {
	// 	fmt.Println("commander:", req.Command, "argser:", req.Args)
	// }

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go handleConn(conn)
	}

}
