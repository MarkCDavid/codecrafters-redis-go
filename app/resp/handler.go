package resp

import (
	"fmt"
	"net"
)

type Handler struct {
	connection net.Conn
	store      *map[string]string
}

func NewHandler(
	connection net.Conn,
	store *map[string]string,
) Handler {
	return Handler{
		connection: connection,
		store:      store,
	}
}

func (handler *Handler) HandleSet(reader *Reader) error {
	key, err := reader.ParseBulkString()
	if err != nil {
		return err
	}

	value, err := reader.ParseBulkString()
	if err != nil {
		return err
	}

	fmt.Println((*handler.store))

	(*handler.store)[key] = value

	_, err = handler.connection.Write(EncodeSimpleString("OK"))
	if err != nil {
		return err
	}

	return nil
}

func (handler *Handler) HandleGet(reader *Reader) error {
	key, err := reader.ParseBulkString()
	if err != nil {
		return err
	}

	value, ok := (*handler.store)[key]

	if ok {
		_, err = handler.connection.Write(EncodeBulkString(value))
		if err != nil {
			return err
		}
	} else {
		_, err = handler.connection.Write(EncodeSimpleString("-1"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *Handler) HandlePing(reader *Reader) error {
	_, err := handler.connection.Write(EncodeSimpleString("PONG"))
	if err != nil {
		return err
	}

	return nil
}

func (handler *Handler) HandleEcho(reader *Reader) error {
	value, err := reader.ParseBulkString()
	if err != nil {
		return err
	}

	_, err = handler.connection.Write(EncodeBulkString(value))
	if err != nil {
		return err
	}

	return nil
}
