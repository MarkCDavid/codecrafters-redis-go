package resp

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type Handler struct {
	connection net.Conn
	Store      storage.Store
}

func NewHandler(
	connection net.Conn,
	store storage.Store,
) Handler {
	return Handler{
		connection: connection,
		Store:      store,
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

	var expiresAt *time.Time
	expiresAt = nil

	if expiresInMs != nil {
		expires := time.Now().UTC().Add(time.Duration(*expiresInMs) * time.Millisecond)
		expiresAt = &expires
	}

	handler.Store.Set(key, value, expiresAt)

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

	entry, ok := handler.Store.Get(key)

	if ok {
		_, err = handler.connection.Write(EncodeBulkString(entry.Value))
		if err != nil {
			return err
		}
	} else {
		_, err = handler.connection.Write(EncodeNullBulkString())
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *Handler) HandleConfigGet(reader *Reader) error {
	var entries []string
	for reader.CanRead(1) {
		key, err := reader.ParseBulkString()
		if err != nil {
			return err
		}

		entry, ok := handler.Store.GetConfig(key)
		if !ok {
			_, err = handler.connection.Write(EncodeNullBulkString())
			if err != nil {
				return err
			}
		}

		entries = append(entries, key)
		entries = append(entries, entry.Value)

	}
	_, err := handler.connection.Write(EncodeArray(entries))
	if err != nil {
		return err
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
