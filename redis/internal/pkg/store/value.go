package store

import (
	"bytes"

	"github.com/seetohjinwei/ccfyi/redis/pkg/delay"
)

type Value struct {
	item  Item
	delay *delay.Delay
}

func NewValue(item Item, delay *delay.Delay) *Value {
	ret := &Value{
		item:  item,
		delay: delay,
	}

	return ret
}

// Item returns item, hasExpired.
// The holder of the Value is responsible for taking note that this value has expired.
func (v *Value) Item() (Item, bool) {
	if v == nil {
		return nil, false
	}
	if v.delay.HasExpired() {
		return v.item, false
	}
	return v.item, true
}

func (v *Value) SerialiseExpiry() []byte {
	if v.delay.HasExpired() {
		return nil
	}

	buf := bytes.Buffer{}

	// FD <=> has delay (micro seconds since Unix epoch)
	buf.WriteString("FD")
	buf.Write(v.delay.Serialise())

	return buf.Bytes()
}
