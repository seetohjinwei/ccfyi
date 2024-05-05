package store

type Item interface {
	Serialise() string
	Get() (string, bool)
	Incr() (int64, bool)
	Decr() (int64, bool)
	LPush(strs []string) (int64, bool)
	RPush(strs []string) (int64, bool)
	LRange(start, stop int) ([]string, bool)
	LLen() (int64, bool)
}
