package resp

import (
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type Handler struct {
	connection net.Conn
	store      storage.Store
}

func NewHandler(
	connection net.Conn,
	store storage.Store,
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

	expiresInMs, err := reader.TryParseSetExpiry()
	if err != nil {
		return err
	}

	handler.store.Set(key, value, expiresInMs)

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

	entry, ok := handler.store.Get(key)

	if ok {
		_, err = handler.connection.Write(EncodeBulkString(entry.Value))
		if err != nil {
			return err
		}
	} else {
		_, err = handler.connection.Write(EncodeBulkString("-1"))
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

func (reader *Reader) TryParseSetExpiry() (*int, error) {
	if !reader.CanRead(1) {
		return nil, nil
	}

	startedAt := reader.Index

	parameterName, err := reader.ParseBulkString()
	if err != nil {
		reader.Index = startedAt
		return nil, err
	}

	if strings.ToUpper(parameterName) != CommandSetParameterExpiry {
		reader.Index = startedAt
		return nil, nil
	}

	valueString, err := reader.ParseBulkString()
	if err != nil {
		reader.Index = startedAt
		return nil, err
	}

	value, err := strconv.Atoi(valueString)
	if err != nil {
		reader.Index = startedAt
		return nil, err
	}

	return &value, nil
}
