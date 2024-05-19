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
	if v.delay.HasExpired() || v.delay == nil {
		// if no delay
		return nil
	}

	buf := bytes.Buffer{}

	// FD <=> has delay (micro seconds since Unix epoch)
	buf.WriteString("FD")
	buf.Write(v.delay.Serialise())

	return buf.Bytes()
}

// Equal checks for equality.
// Should only be used for tests.
func (v *Value) Equal(other any) bool {
	o, ok := other.(*Value)
	if !ok {
		return false
	}

	if v == nil || o == nil {
		return (v == nil) && (o == nil)
	}

	return v.delay.Equal(o.delay) && v.item.Equal(o.item)
}
