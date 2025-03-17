package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	BIND_ADDRESS = "0.0.0.0"
	BIND_PORT    = "6379"
)

func main() {
	address := fmt.Sprintf("%s:%s", BIND_ADDRESS, BIND_PORT)

	fmt.Printf("Listening on %s...\n", address)
	socket, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to bind to port %s.\n", BIND_PORT)
		os.Exit(1)
	}

	connection, err := socket.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1000)

	count, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	commands := string(buffer[:count])
	for range strings.Split(commands, "\n") {
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
