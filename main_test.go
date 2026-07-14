package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func newTestConn(t *testing.T) (net.Conn, *SafeMap) {
	t.Helper()
	server, client := net.Pipe()
	store := NewSafeMap()
	go handleConn(server, store)
	t.Cleanup(func() { client.Close() })
	return client, store
}

func send(t *testing.T, conn net.Conn, msg string) string {
	t.Helper()
	fmt.Fprintln(conn, msg)
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	return line
}

func TestHandler_Ping(t *testing.T) {
	conn, _ := newTestConn(t)
	got := send(t, conn, "PING")
	if got != "OK|PONG\n" {
		t.Errorf("PING = %q, want %q", got, "OK|PONG\n")
	}
}

func TestHandler_SetAndGet(t *testing.T) {
	conn, _ := newTestConn(t)
	send(t, conn, "SET|city|accra")
	got := send(t, conn, "GET|city")
	if got != "OK|accra\n" {
		t.Errorf("GET after SET = %q, want %q", got, "OK|accra\n")
	}
}

func TestHandler_GetKeyNotFound(t *testing.T) {
	conn, _ := newTestConn(t)
	got := send(t, conn, "GET|missing")
	if got != "ERR|key not found\n" {
		t.Errorf("GET missing = %q, want %q", got, "ERR|key not found\n")
	}
}

func TestHandler_DelExistingKey(t *testing.T) {
	conn, _ := newTestConn(t)
	send(t, conn, "SET|x|1")
	got := send(t, conn, "DEL|x")
	if got != "OK|\n" {
		t.Errorf("DEL = %q, want %q", got, "OK|\n")
	}
}

func TestHandler_DelKeyNotFound(t *testing.T) {
	conn, _ := newTestConn(t)
	got := send(t, conn, "DEL|missing")
	if got != "ERR|key not found\n" {
		t.Errorf("DEL missing = %q, want %q", got, "ERR|key not found\n")
	}
}

func TestHandler_List(t *testing.T) {
	conn, _ := newTestConn(t)
	send(t, conn, "SET|a|1")
	send(t, conn, "SET|b|2")
	got := send(t, conn, "LIST")
	if got != "OK|a|b\n" && got != "OK|b|a\n" {
		t.Errorf("LIST = %q, want OK|a|b or OK|b|a", got)
	}
}

func TestHandler_InvalidCommand(t *testing.T) {
	conn, _ := newTestConn(t)
	got := send(t, conn, "BADCMD")
	if got[:3] != "ERR" {
		t.Errorf("invalid command response = %q, want ERR prefix", got)
	}
}
