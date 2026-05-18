package main

import (
	"bufio"
	"fmt"
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

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				line := scanner.Text()
				cmd, err := parser.ParseCommand(line)
				if err != nil {
					log.Println("parse error:", err)
					continue
				}
				fmt.Println("got command:", cmd)
			}

			log.Println("Connection closed", conn.RemoteAddr())
		}(conn)
	}
}
