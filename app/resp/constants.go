package resp

const (
	RespSimpleStringType = "+"
	RespArrayType        = "*"
	RespBulkStringType   = "$"
	RespEOL              = "\r\n"

	CommandEcho = "ECHO"
	CommandPing = "PING"

	CommandConfig = "CONFIG"

	CommandSet                = "SET"
	CommandSetParameterExpiry = "PX"

	CommandGet = "GET"

	CommandKeys = "KEYS"
)
