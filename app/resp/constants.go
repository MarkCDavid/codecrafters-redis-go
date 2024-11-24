package resp

const (
	RespSimpleStringType = "+"
	RespArrayType        = "*"
	RespBulkStringType   = "$"
	RespEOL              = "\r\n"

	CommandEcho = "ECHO"
	CommandPing = "PING"
	CommandSet  = "SET"
	CommandGet  = "GET"
)
