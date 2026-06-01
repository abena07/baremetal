package main

import (
	"fmt"
	"slices"
	"strings"
)

type Request struct {
	Command string
	Args    []string
}

func ParseRequest(input string) (Request, error) {
	commands := []string{"SET", "GET", "PING"}

	formatted_input := strings.TrimRight(input, "\r\n")
	line := strings.Split(formatted_input, "|")

	fmt.Println("args", line)

	if len(line) < 1 {
		return Request{}, fmt.Errorf("empty message")
	}

	command := line[0]
	args := line[1:]

	if !slices.Contains(commands, command) {
		return Request{}, fmt.Errorf("unknown command: %q", command)
	}

	return Request{Command: command, Args: args}, nil

}

func ResponseOkFormatter(command string, values []string) string {
	success := "OK"

	if command == "PING" {
		return success + "|" + "PONG" + "\n"
	}

	if command == "GET" {
		return success + "|" + values[1] + "\n"
	}

	return success + "|" + "\n"
}

func ResponseErrFormatter(message string) string {
	fail := "ERR"

	return fail + "|" + message + "\n"

}
