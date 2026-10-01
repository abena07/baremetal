package main

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	Command string
	Args    []string
}

func ParseRequest(input string) (Request, error) {
	input = strings.TrimSuffix(input, "\r")
	if strings.TrimSpace(input) == "" {
		return Request{}, fmt.Errorf("empty message")
	}

	name, rest, hasArgs := strings.Cut(input, "|")
	command := strings.TrimSpace(name)

	spec, ok := commands[command]
	if !ok {
		return Request{}, fmt.Errorf("unknown command: %q", command)
	}

	args := []string{}
	if hasArgs {
		if spec.TrailingValue {
			args = strings.SplitN(rest, "|", spec.Arity)
		} else {
			args = strings.Split(rest, "|")
		}
	}

	if len(args) != spec.Arity {
		return Request{}, arityError(command, spec.Arity)
	}

	if spec.Arity > 0 {
		args[0] = strings.TrimSpace(args[0])
		if args[0] == "" {
			return Request{}, fmt.Errorf("key must not be empty")
		}
	}

	return Request{Command: command, Args: args}, nil
}

func arityError(command string, arity int) error {
	switch arity {
	case 0:
		return fmt.Errorf("%s requires 0 arguments", command)
	case 1:
		return fmt.Errorf("%s requires exactly 1 argument", command)
	default:
		return fmt.Errorf("%s requires exactly %d arguments", command, arity)
	}
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
