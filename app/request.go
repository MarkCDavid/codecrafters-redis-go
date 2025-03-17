package main

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	SPECIFIER_ARRAY       = byte('*')
	SPECIFIER_BULK_STRING = byte('$')
)

var (
	TERMINATOR = []byte{13, 10}
)

func ParseCommand(value []byte) ([]*string, error) {
	if err := ensureStartsWith(value, SPECIFIER_ARRAY); err != nil {
		return nil, err
	}
	value = value[1:]

	elementCount, value, err := parseLength(value)
	if err != nil {
		return nil, err
	}

	result := make([]*string, 0, elementCount)
	for i := int64(0); i < elementCount; i++ {
		if err := ensureStartsWith(value, SPECIFIER_BULK_STRING); err != nil {
			return nil, err
		}

		stringLength, value, err := parseLength(value[1:])
		if err != nil {
			return nil, err
		}

		if stringLength == -1 {
			result = append(result, nil)
			continue
		}

		stringValue, value, err := parseString(value, stringLength)
		if err != nil {
			return nil, err
		}

		result = append(result, stringValue)
	}

	return result, nil
}

func parseString(value []byte, length int64) (*string, []byte, error) {
	encodedString, leftoverValue, found := bytes.Cut(value, TERMINATOR)
	if !found {
		return nil, value, fmt.Errorf("Terminator could not be found.")
	}

	if len(encodedString) != int(length) {
		return nil, value, fmt.Errorf("Provided string is larger (%d) than expected (%d).", len(encodedString), length)
	}

	result := string(encodedString)
	return &result, leftoverValue, nil
}
func parseLength(value []byte) (int64, []byte, error) {
	encodedLength, leftoverValue, found := bytes.Cut(value, TERMINATOR)
	if !found {
		return -1, value, fmt.Errorf("Terminator could not be found.")
	}

	length, err := strconv.ParseInt(string(encodedLength), 10, 0)
	if err != nil {
		return -1, value, fmt.Errorf(
			"Failed to convert '%s' into an integer.",
			string(encodedLength),
		)
	}

	return length, leftoverValue, nil
}

func ensureStartsWith(haystack []byte, needle byte) error {
	if haystack[0] != needle {
		return fmt.Errorf(
			"Expected '%c' as the first command byte, got '%c'.",
			rune(needle),
			rune(haystack[0]),
		)
	}
	return nil
}
