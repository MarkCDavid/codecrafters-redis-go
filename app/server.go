package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/app/connection"
)

const (
	pingResponse = "+PONG\r\n"
)

func registerInterruptHandling(cancel context.CancelFunc) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChannel
		fmt.Println("Received interrupt. Signaling cancellation.")
		cancel()
	}()
}

func handleConnection(
	ctx context.Context,
	connection net.Conn,
) {
	defer connection.Close()
	for {
		select {
		case <-ctx.Done():
			return

		default:
			command := make([]byte, 1024)
			_, err := connection.Read(command)

			if errors.Is(err, io.EOF) {
				fmt.Println("Client closed the connecion.")
				return
			}

			if err != nil {
				fmt.Println("Error reading command:", err.Error())
				return
			}

			_, err = connection.Write([]byte(pingResponse))
			if err != nil {
				fmt.Println("Error writing response:", err.Error())
				return
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	registerInterruptHandling(cancel)

	socket := connection.StartServer(ctx, "6379")
	var wg sync.WaitGroup

	connection.AcceptConnections(ctx, &wg, socket, handleConnection)
	fmt.Println("Waiting for all connections to close.")
	wg.Wait()
	fmt.Println("All connections closed. Exiting...")
}
