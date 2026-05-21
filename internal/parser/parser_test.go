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

func TestUnSubMissingTopic(t *testing.T) {
	cmd, err := ParseCommand(makeReader("UNSUB \n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestUnSubMissingArgs(t *testing.T) {
	cmd, err := ParseCommand(makeReader("UNSUB\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestUnSubInvalidSid(t *testing.T) {
	cmd, err := ParseCommand(makeReader("UNSUB notanumber\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for non-numeric sid, got nil")
	}
}

func TestUnSubHappyPath(t *testing.T) {
	cmd, err := ParseCommand(makeReader("UNSUB 42\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unsub, ok := cmd.(UnsubCommand)
	if !ok {
		t.Fatalf("expected UnsubCommand, got %T", cmd)
	}
	if unsub.Sid != 42 {
		t.Errorf("got sid %d, want %d", unsub.Sid, 42)
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

func TestPubMissingArgs(t *testing.T) {
	cmd, err := ParseCommand(makeReader("PUB\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestPubMissingByteCount(t *testing.T) {
	cmd, err := ParseCommand(makeReader("PUB sports.scores\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestPubInvalidByteCount(t *testing.T) {
	cmd, err := ParseCommand(makeReader("PUB sports.scores notanumber\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for non-numeric byte count, got nil")
	}
}

func TestPubNegativeByteCount(t *testing.T) {
	cmd, err := ParseCommand(makeReader("PUB sports.scores -1\r\n\r\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for negative byte count, got nil")
	}
}

func TestPubExceedsMaxPayload(t *testing.T) {
	cmd, err := ParseCommand(makeReader("PUB sports.scores 1048577\r\n\r\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for oversized byte count, got nil")
	}
}

func TestPubPayloadShorterThanDeclared(t *testing.T) {
	// declares 10 bytes but only provides 3
	cmd, err := ParseCommand(makeReader("PUB sports.scores 10\r\nHi\r\n"))
	if cmd != nil {
		t.Errorf("expected nil command, got %v", cmd)
	}
	if err == nil {
		t.Error("expected an error for truncated payload, got nil")
	}
}

func TestPubHappyPath(t *testing.T) {
	cmd, err := ParseCommand(makeReader("PUB sports.scores 5\r\nHello\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pub, ok := cmd.(PubCommand)
	if !ok {
		t.Fatalf("expected PubCommand, got %T", cmd)
	}
	if pub.Topic != "sports.scores" {
		t.Errorf("got topic %q, want %q", pub.Topic, "sports.scores")
	}
	if pub.Payload != "Hello" {
		t.Errorf("got payload %q, want %q", pub.Payload, "Hello")
	}
}

// Verifies the reader is positioned correctly after a PUB so the next
// command can be parsed without corruption.
func TestPubFollowedByNextCommand(t *testing.T) {
	r := makeReader("PUB sports.scores 5\r\nHello\r\nSUB sports.scores 1\r\n")

	_, err := ParseCommand(r)
	if err != nil {
		t.Fatalf("first ParseCommand failed: %v", err)
	}

	cmd, err := ParseCommand(r)
	if err != nil {
		t.Fatalf("second ParseCommand failed: %v", err)
	}
	if _, ok := cmd.(SubCommand); !ok {
		t.Fatalf("expected SubCommand after PUB, got %T", cmd)
	}
}
