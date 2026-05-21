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
					} else {
						log.Println("parse error:", err)
					}
					return
				}
				fmt.Println("got command:", cmd)
			}

		}(conn)
	}
}
