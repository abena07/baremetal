package main

import (
	"bytes"
	"testing"
)

func TestParseRequest_ValidCommands(t *testing.T) {
	tests := []struct {
		input   string
		command string
		args    []string
	}{
		{"SET|key|value", "SET", []string{"key", "value"}},
		{"GET|key", "GET", []string{"key"}},
		{"DEL|key", "DEL", []string{"key"}},
		{"PING", "PING", []string{}},
		{"LIST", "LIST", []string{}},
	}

	for _, tt := range tests {
		req, err := ParseRequest(tt.input)
		if err != nil {
			t.Errorf("ParseRequest(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if req.Command != tt.command {
			t.Errorf("ParseRequest(%q) command = %q, want %q", tt.input, req.Command, tt.command)
		}
		if len(req.Args) != len(tt.args) {
			t.Errorf("ParseRequest(%q) args = %v, want %v", tt.input, req.Args, tt.args)
		}
	}
}

func TestParseRequest_Errors(t *testing.T) {
	tests := []struct {
		input string
		desc  string
	}{
		{"", "empty input"},
		{"UNKNOWN|arg", "unknown command"},
		{"GET", "GET missing arg"},
		{"DEL", "DEL missing arg"},
		{"SET|key", "SET missing value"},
		{"GET|key|extra", "GET too many args"},
		{"PING|extra", "PING with args"},
		{"LIST|extra", "LIST with args"},
		{"SET|key|val|extra|extra2", "too many pipes"},
	}

	for _, tt := range tests {
		_, err := ParseRequest(tt.input)
		if err == nil {
			t.Errorf("ParseRequest(%q) (%s) expected error, got nil", tt.input, tt.desc)
		}
	}
}

func TestWriteOK(t *testing.T) {
	var buf bytes.Buffer
	WriteOK(&buf, "PONG")
	if got := buf.String(); got != "OK|PONG\n" {
		t.Errorf("WriteOK = %q, want %q", got, "OK|PONG\n")
	}
}

func TestWriteOK_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	WriteOK(&buf, "")
	if got := buf.String(); got != "OK|\n" {
		t.Errorf("WriteOK empty = %q, want %q", got, "OK|\n")
	}
}

func TestWriteErr(t *testing.T) {
	var buf bytes.Buffer
	WriteErr(&buf, "key not found")
	if got := buf.String(); got != "ERR|key not found\n" {
		t.Errorf("WriteErr = %q, want %q", got, "ERR|key not found\n")
	}
}
