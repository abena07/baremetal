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
	commands := []string{"SET", "GET", "PING", "DEL"}

	formatted_input := strings.TrimRight(input, "\r\n")

	if strings.Count(formatted_input, "|") > 2 {
		return Request{}, fmt.Errorf("invalid character in value")
	}

	raw_line := strings.Split(formatted_input, "|")

	var line []string
	for _, item := range raw_line {
		cleaned := strings.TrimSpace(item)
		if cleaned != "" {
			line = append(line, cleaned)
		}
	}

	if len(line) < 1 {
		return Request{}, fmt.Errorf("empty message")
	}

	command := line[0]
	args := line[1:]

	if !slices.Contains(commands, command) {
		return Request{}, fmt.Errorf("unknown command: %q", command)
	}

	if command == "PING" {
		if len(args) != 0 {
			return Request{}, fmt.Errorf("%s requires 0 arguments", command)
		}
	} else {
		if len(args) != 2 {
			return Request{}, fmt.Errorf("%s requires exactly 2 arguments", command)
		}
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
