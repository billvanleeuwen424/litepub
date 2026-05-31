package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"litepub/internal/parser"
)

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

				// Goal for today
				// fanout on reciept of message.
				// the parser just parsed and returned us a struct containing the message
				// we can get a SUB, PUB, ACK, or UNSUB
				// Functions need to be called for each

				log.Println("got command:", cmd)

				switch cmd.(type) {
				case parser.SubCommand, parser.PubCommand, parser.UnsubCommand, parser.AckCommand:
					if _, err = fmt.Fprintf(conn, "+OK\r\n"); err != nil {
						log.Println("write error:", err)
						return
					}
				}
			}

		}(conn)
	}
}
