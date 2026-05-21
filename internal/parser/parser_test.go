package parser

import (
	"bufio"
	"strings"
	"testing"
)

func makeReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

// My first ever Go unit test :)
func TestSubMissingArgs(t *testing.T) {
	cmd, err := ParseCommand(makeReader("SUB\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestSubMissingTopic(t *testing.T) {
	cmd, err := ParseCommand(makeReader("SUB \n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestSubHappyPath(t *testing.T) {
	cmd, err := ParseCommand(makeReader("SUB sports.scores 42\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sub, ok := cmd.(SubCommand)
	if !ok {
		t.Fatalf("expected SubCommand, got %T", cmd)
	}
	if sub.Topic != "sports.scores" {
		t.Errorf("got topic %q, want %q", sub.Topic, "sports.scores")
	}
	if sub.Sid != 42 {
		t.Errorf("got sid %d, want %d", sub.Sid, 42)
	}
}

func TestSubInvalidSid(t *testing.T) {
	cmd, err := ParseCommand(makeReader("SUB sports.scores notanumber\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for non-numeric sid, got nil")
	}
}

func TestEmptyInput(t *testing.T) {
	cmd, err := ParseCommand(makeReader(""))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for empty input, got nil")
	}
}

func TestUnknownCommand(t *testing.T) {
	cmd, err := ParseCommand(makeReader("BLAH foo 1\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for unknown command, got nil")
	}
}
