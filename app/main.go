package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	BIND_ADDRESS = "0.0.0.0"
	BIND_PORT    = "6379"
)

func main() {
	address := fmt.Sprintf("%s:%s", BIND_ADDRESS, BIND_PORT)

	fmt.Printf("Listening on %s...\n", address)
	socket, err := net.Listen("tcp", address)
	defer socket.Close()
	if err != nil {
		fmt.Printf("Failed to bind to port %s.\n", BIND_PORT)
		os.Exit(1)
	}

	for {
		connection, err := socket.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}

		go HandleConnection(connection)
	}
}

func HandleConnection(connection net.Conn) {
	defer connection.Close()
	for {
		buffer := make([]byte, 100)
		_, err := connection.Read(buffer)

		if errors.Is(err, io.EOF) {
			fmt.Println("Connection closed.")
			return
		}
		if err != nil {
			fmt.Println("Error reading data: ", err.Error())
			return
		}

		connection.Write(AsSimpleString("PONG"))
	}
}

// https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings
func AsSimpleString(value string) []byte {
	simpleString := make([]byte, len(value)+3)
	simpleString[0] = byte('+')
	for index, character := range value {
		simpleString[index+1] = byte(character)
	}

	simpleString[len(simpleString)-2] = byte('\r')
	simpleString[len(simpleString)-1] = byte('\n')

	return simpleString
}
