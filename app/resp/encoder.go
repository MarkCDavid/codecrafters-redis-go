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
