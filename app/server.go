package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	pingResponse = "+PONG\r\n"
)

func main() {
	socket, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	connection, err := socket.Accept()
	defer connection.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {
		command := make([]byte, 1024)
		_, err = connection.Read(command)

		if errors.Is(err, io.EOF) {
			fmt.Println("Client closed the connecion.")
			return
		}

		if err != nil {
			fmt.Println("Error reading command:", err.Error())
			os.Exit(1)
		}

		_, err = connection.Write([]byte(pingResponse))
		if err != nil {
			fmt.Println("Error writing response:", err.Error())
			os.Exit(1)
		}
	}
}
