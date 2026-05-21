package parser

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type SubCommand struct {
	Sid   int
	Topic string
}

// type PubCommand struct {
// 	Bytes   int
// 	Topic   string
// 	Payload string
// }

type Command interface {
	CommandType()
}

func (s SubCommand) CommandType() {}

func ParseCommand(r *bufio.Reader) (Command, error) {

	command, err := r.ReadString('\n')

	if err != nil {
		return nil, err
	}

	words := strings.Fields(command)

	if len(words) == 0 {
		return nil, fmt.Errorf("bad input")
	}

	if words[0] == "SUB" {

		if len(words) < 3 {
			return nil, fmt.Errorf("bad input")
		}

		subid, err := strconv.Atoi(words[2])
		if err != nil {
			return nil, err
		}

		return SubCommand{subid, words[1]}, nil
	}

	return nil, fmt.Errorf("unknown command, not implemented: %s", words[0])

}
