package delay

import (
	"time"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
)

type Delay struct {
	expiry time.Time
}

func NewDelay(expiry time.Time) *Delay {
	ret := &Delay{
		expiry: expiry,
	}
	return ret
}

func (d *Delay) HasExpired() bool {
	if d == nil {
		return false
	}

	return time.Now().After(d.expiry)
}

func (d *Delay) Serialise() []byte {
	return encoding.EncodeInteger(d.expiry.UnixMicro())
}
