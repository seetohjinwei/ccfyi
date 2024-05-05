package store

type Item interface {
	Get() (string, bool)
	Incr() (int64, bool)
	Decr() (int64, bool)
}
