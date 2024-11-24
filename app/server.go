package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/app/connection"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
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
	store storage.Store,
) {
	defer connection.Close()
	for {
		select {
		case <-ctx.Done():
			return

		default:
			readBuffer := make([]byte, 1024)
			readLength, err := connection.Read(readBuffer)

			if errors.Is(err, io.EOF) {
				fmt.Println("Client closed the connecion.")
				return
			}

			if err != nil {
				fmt.Println("Error reading command:", err.Error())
				return
			}

			handler := resp.NewHandler(connection, store)
			err = resp.Parse(&handler, readBuffer, readLength)
			if err != nil {
				fmt.Println("Failed to parse command:", err.Error())
				continue
			}

			// _, err = connection.Write([]byte(pingResponse))
			// if err != nil {
			// 	fmt.Println("Error writing response:", err.Error())
			// 	return
			// }
		}
	}
}

func main() {
	rdbDirectory := flag.String("dir", "/tmp/rdb", "")
	rdbFileName := flag.String("dbfilename", "dump.rdb", "")

	store := storage.NewStore()

	store.SetConfig("dir", *rdbDirectory, nil)
	store.SetConfig("dbfilename", *rdbFileName, nil)

	ctx, cancel := context.WithCancel(context.Background())
	registerInterruptHandling(cancel)

	socket := connection.StartServer(ctx, "6379")
	var wg sync.WaitGroup

	connection.AcceptConnections(ctx, &wg, socket, store, handleConnection)
	fmt.Println("Waiting for all connections to close.")
	wg.Wait()
	fmt.Println("All connections closed. Exiting...")
}
