package resp

import "net"

type Handler struct {
	connection net.Conn
}

func NewHandler(
	connection net.Conn,
) Handler {
	return Handler{
		connection: connection,
	}
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
