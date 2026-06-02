package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	Command string
	Args    []string
}

func ParseRequest(input string) (Request, error) {
	commands := []string{"SET", "GET", "PING", "DEL"}
	oneArgCommands := []string{"GET", "DEL"}

	if strings.Count(input, "|") > 2 {
		return Request{}, fmt.Errorf("invalid character in value")
	}

	rawLine := strings.Split(input, "|")

	var line []string
	for _, item := range rawLine {
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
	} else if slices.Contains(oneArgCommands, command) {
		if len(args) != 1 {
			return Request{}, fmt.Errorf("%s requires exactly 1 argument", command)
		}
	} else {
		if len(args) != 2 {
			return Request{}, fmt.Errorf("%s requires exactly 2 arguments", command)
		}
	}

	return Request{Command: command, Args: args}, nil
}

func WriteOK(w io.Writer, result string) {
	okPrefix := "OK"
	response := okPrefix + "|" + result + "\n"
	w.Write([]byte(response))
}

func WriteErr(w io.Writer, message string) {
	errPrefix := "ERR"
	response := errPrefix + "|" + message + "\n"
	w.Write([]byte(response))
}
