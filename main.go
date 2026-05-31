package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"litepub/internal/parser"
)

const okResponse = "+OK\r\n"

func writeOK(conn net.Conn) error {
	if _, err := io.WriteString(conn, okResponse); err != nil {
		return err
	}

	return nil
}

func handleSubCommand(conn net.Conn, s parser.SubCommand) error {
	log.Println("got SUB command:", s)

	return writeOK(conn)
}

func handlePubCommand(conn net.Conn, s parser.PubCommand) error {
	log.Println("got PUB command:", s)

	return writeOK(conn)
}

func handleAckCommand(conn net.Conn, s parser.AckCommand) error {
	log.Println("got ACK command:", s)

	return writeOK(conn)
}

func handleUnsubCommand(conn net.Conn, s parser.UnsubCommand) error {
	log.Println("got UNSUB command:", s)

	return writeOK(conn)
}

func main() {
	fmt.Println("Starting server")

	// Flush coverage data and exit cleanly when signalled.
	// go build -cover registers counter-flush as an os.Exit hook, but
	// SIGTERM/SIGINT bypass that hook entirely. We catch both signals here
	// so that integration-test suites that terminate the broker via
	// proc.terminate() (SIGTERM) still get accurate coverage output.
	go handleSignals()

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		go func(conn net.Conn) {
			log.Println("Connection from", conn.RemoteAddr())
			defer conn.Close()

			reader := bufio.NewReader(conn)
			for {
				cmd, err := parser.ParseCommand(reader)
				if err != nil {
					if err == io.EOF {
						log.Println("client disconnected", conn.RemoteAddr())
						return
					}
					log.Println("parse error:", err)
					if _, err = fmt.Fprintf(conn, "-ERR %s\r\n", err); err != nil {
						log.Println("write error:", err)
						return
					}
					continue
				}

				switch c := cmd.(type) {
				case parser.SubCommand:
					err = handleSubCommand(conn, c)
				case parser.PubCommand:
					err = handlePubCommand(conn, c)
				case parser.UnsubCommand:
					err = handleUnsubCommand(conn, c)
				case parser.AckCommand:
					err = handleAckCommand(conn, c)
				default:
					log.Println("unhandled command type:", c)
				}

				if err != nil {
					log.Println(err)
					return
				}
			}

		}(conn)
	}
}
