package delay

import "time"

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
