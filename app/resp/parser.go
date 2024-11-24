package resp

import (
	"fmt"
	"strings"
)

func Parse(
	handler *Handler,
	buffer []byte,
	size int,
) error {
	reader := NewReader(buffer, size)

	_, err := reader.ParseArray()
	if err != nil {
		return err
	}

	for reader.CanRead(1) {
		value, err := reader.ParseBulkString()
		if err != nil {
			return err
		}

		switch strings.ToUpper(value) {
		case CommandPing:
			err = handler.HandlePing(&reader)
		case CommandEcho:
			err = handler.HandleEcho(&reader)
		case CommandSet:
			err = handler.HandleSet(&reader)
		case CommandGet:
			err = handler.HandleGet(&reader)
		case CommandConfig:
			err = ParseConfig(handler, &reader)
		default:
			fmt.Println(string(buffer))
			return fmt.Errorf("Command %s is not implemented.", value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func ParseConfig(
	handler *Handler,
	reader *Reader,
) error {
	value, err := reader.ParseBulkString()
	if err != nil {
		return err
	}

	switch strings.ToUpper(value) {
	case CommandGet:
		err = handler.HandleConfigGet(reader)
	default:
		return fmt.Errorf("Command %s %s is not implemented.", CommandConfig, value)
	}

	if err != nil {
		return err
	}

	return nil
}

func (reader *Reader) ParseBulkString() (string, error) {
	err := reader.IsExpectedType(RespBulkStringType)
	if err != nil {
		return "", err
	}

	length, err := reader.ReadInteger()
	if err != nil {
		return "", err
	}

	err = reader.ParseEOL()
	if err != nil {
		return "", err
	}

	value, err := reader.Read(length)
	if err != nil {
		return "", err
	}

	err = reader.ParseEOL()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (reader *Reader) ParseArray() (int, error) {
	err := reader.IsExpectedType(RespArrayType)
	if err != nil {
		return 0, err
	}

	length, err := reader.ReadInteger()
	if err != nil {
		return 0, err
	}

	err = reader.ParseEOL()
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (reader *Reader) IsExpectedType(expectedRespType string) error {
	actualRespType, err := reader.Read(1)
	if err != nil {
		return err
	}

	if expectedRespType != actualRespType {
		return fmt.Errorf("Client sent invalid command. Expected %s, got %s.", expectedRespType, actualRespType)
	}

	return nil
}

func (reader *Reader) ParseEOL() error {
	_, err := reader.Read(2)
	if err != nil {
		return err
	}
	return nil
}
