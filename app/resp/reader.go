package resp

import (
	"fmt"
	"strconv"
)

type Reader struct {
	buffer []byte
	size   int
	index  int
}

func NewReader(
	buffer []byte,
	size int,
) Reader {
	return Reader{
		buffer: buffer,
		size:   size,
		index:  0,
	}
}

func (reader *Reader) ReadInteger() (int, error) {
	result := ""
	for {
		value, err := reader.Peek(1)
		if err != nil {
			return 0, err
		}

		if !isDigit(value) {
			resultInteger, err := strconv.Atoi(result)
			if err != nil {
				return 0, nil
			}

			return resultInteger, nil
		}

		reader.index += 1
		result += value
	}
}

func (reader *Reader) Read(
	length int,
) (string, error) {
	value, err := reader.Peek(length)
	if err != nil {
		return "", err
	}

	reader.index += length
	return value, nil
}

func (reader *Reader) Peek(
	length int,
) (string, error) {
	if !reader.CanRead(length) {
		return "", fmt.Errorf("%d exceeds size of %d", reader.index+length, reader.size)
	}

	value := string(reader.buffer[reader.index : reader.index+length])
	return value, nil
}

func (reader *Reader) CanRead(
	length int,
) bool {
	return reader.index+length <= reader.size
}

func isDigit(value string) bool {
	if len(value) != 1 {
		return false
	}

	return value[0] >= '0' && value[0] <= '9'
}
