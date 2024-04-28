package messages

const (
	CR   byte   = '\r'
	LF   byte   = '\n'
	CRLF string = "\r\n"
)

// https://redis.io/docs/latest/develop/reference/protocol-spec/
type Message interface {
	Serialise() string
}
