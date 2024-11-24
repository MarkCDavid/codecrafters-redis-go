package connection

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func StartServer(
	ctx context.Context,
	port string,
) net.Listener {
	socket, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", port)
		os.Exit(1)
	}

	fmt.Printf("Listening on port %s\n", port)
	go func() {
		<-ctx.Done()
		fmt.Println("Received interrupt. Closing listener.")
		socket.Close()
	}()

	return socket
}

func AcceptConnections(
	ctx context.Context,
	wg *sync.WaitGroup,
	socket net.Listener,
	store storage.Store,
	handler func(context.Context, net.Conn, storage.Store),
) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Received interrupt. Stopping connection acceptance.")
			return
		default:
			connection := acceptConnection(socket)
			if connection == nil {
				continue
			}

			fmt.Printf("Received connection from %s\n", connection.RemoteAddr().String())
			go func() {
				wg.Add(1)
				handler(ctx, connection, store)
				wg.Done()
			}()
		}
	}
}

func acceptConnection(
	socket net.Listener,
) net.Conn {
	connection, err := socket.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		return nil
	}

	return connection
}
