package main

import (
	"errors"
	"strings"
)

var errKeyNotFound = errors.New("key not found")

type CommandSpec struct {
	// Arity is the exact number of arguments the command takes.
	Arity int
	// TrailingValue lets the last argument contain `|`, so values aren't restricted.
	TrailingValue bool
	Handler       func(s *SafeMap, args []string) (string, error)
}

// commands is the single source of truth for what the server understands.
// adding a command means adding an entry here; the parser and dispatcher pick it up.
var commands = map[string]CommandSpec{
	"PING": {
		Arity: 0,
		Handler: func(_ *SafeMap, _ []string) (string, error) {
			return "PONG", nil
		},
	},
	"SET": {
		Arity:         2,
		TrailingValue: true,
		Handler: func(s *SafeMap, args []string) (string, error) {
			s.Set(args[0], args[1])
			return "", nil
		},
	},
	"GET": {
		Arity: 1,
		Handler: func(s *SafeMap, args []string) (string, error) {
			val, ok := s.Get(args[0])
			if !ok {
				return "", errKeyNotFound
			}
			return val, nil
		},
	},
	"DEL": {
		Arity: 1,
		Handler: func(s *SafeMap, args []string) (string, error) {
			if !s.Delete(args[0]) {
				return "", errKeyNotFound
			}
			return "", nil
		},
	},
	"LIST": {
		Arity: 0,
		Handler: func(s *SafeMap, _ []string) (string, error) {
			return strings.Join(s.List(), "|"), nil
		},
	},
}
