package messages

// A null bulk string is represented by `nil`.
type BulkString struct {
	len uint
	str string
}

func (r *BulkString) Serialise() string {
	return r.str
}
