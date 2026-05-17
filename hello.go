package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("Strarting server")

	listener, err := net.Listen("tcp", ":8080")

	// check for errors on listening
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
			io.Copy(conn, conn)
			conn.Close()
			log.Println("Connection closed", conn.RemoteAddr())
		}(conn)

	}
}
