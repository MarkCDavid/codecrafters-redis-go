package resp

import "fmt"

func EncodeSimpleString(value string) []byte {
	return []byte(fmt.Sprintf(
		"%s%s%s",
		RespSimpleStringType,
		value,
		RespEOL,
	))
}

func EncodeNullBulkString() []byte {
	return []byte(fmt.Sprintf(
		"%s%d%s",
		RespBulkStringType,
		-1,
		RespEOL,
	))
}

func EncodeBulkString(value string) []byte {
	return []byte(fmt.Sprintf(
		"%s%d%s%s%s",
		RespBulkStringType,
		len(value),
		RespEOL,
		value,
		RespEOL,
	))
}

func EncodeArray(values []string) []byte {
	result := fmt.Sprintf(
		"%s%d%s",
		RespArrayType,
		len(values),
		RespEOL,
	)

	for _, value := range values {
		result += string(EncodeBulkString(value))
	}

	return []byte(result)
}
