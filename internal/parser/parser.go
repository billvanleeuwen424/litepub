package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type SubCommand struct {
	Sid   int
	Topic string
}

type UnsubCommand struct {
	Sid int
}

type AckCommand struct {
	MsgId int
}

type PubCommand struct {
	Topic   string
	Payload string
}

const maxPayloadBytes = 1 << 20 // 1 MB

type Command interface {
	CommandType()
}

func (s SubCommand) CommandType() {}
func (s SubCommand) String() string {
	return fmt.Sprintf("SUB topic=%s sid=%d", s.Topic, s.Sid)
}

func (s PubCommand) CommandType() {}
func (s PubCommand) String() string {
	return fmt.Sprintf("PUB topic=%s payload=%s", s.Topic, s.Payload)
}

func (s UnsubCommand) CommandType() {}
func (s UnsubCommand) String() string {
	return fmt.Sprintf("UNSUB sid=%d", s.Sid)
}

func (s AckCommand) CommandType() {}
func (s AckCommand) String() string {
	return fmt.Sprintf("ACK MsgId=%d", s.MsgId)
}

func ParseCommand(r *bufio.Reader) (Command, error) {

	command, err := r.ReadString('\n')

	if err != nil {
		return nil, err
	}

	words := strings.Fields(command)

	if len(words) == 0 {
		return nil, fmt.Errorf("bad input")
	}

	switch words[0] {
	case "SUB":

		if len(words) < 3 {
			return nil, fmt.Errorf("bad input")
		}

		subid, err := strconv.Atoi(words[2])
		if err != nil {
			return nil, err
		}

		return SubCommand{subid, words[1]}, nil
	case "UNSUB":

		if len(words) < 2 {
			return nil, fmt.Errorf("bad input")
		}

		subid, err := strconv.Atoi(words[1])
		if err != nil {
			return nil, err
		}

		return UnsubCommand{subid}, nil
	case "ACK":
		if len(words) < 2 {
			return nil, fmt.Errorf("bad input")
		}

		msgId, err := strconv.Atoi(words[1])
		if err != nil {
			return nil, err
		}

		return AckCommand{msgId}, nil
	case "PUB":
		if len(words) < 3 {
			return nil, fmt.Errorf("bad input")
		}

		bytes, err := strconv.Atoi(words[2])
		if err != nil {
			return nil, err
		}

		if bytes < 0 || bytes > maxPayloadBytes {
			return nil, fmt.Errorf("bad input")
		}

		buf := make([]byte, bytes)
		_, err = io.ReadFull(r, buf)

		if err != nil {
			return nil, err
		}

		_, err = r.Discard(2) // consume trailing \r\n after payload
		if err != nil {
			return nil, err
		}

		return PubCommand{words[1], string(buf)}, nil
	}

	return nil, fmt.Errorf("unknown command, not implemented: %s", words[0])

}
