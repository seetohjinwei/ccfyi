package store

type AbstractItem struct{}

func (b *AbstractItem) Get() (string, bool) {
	return "", false
}

func (b *AbstractItem) Incr() (int64, bool) {
	return 0, false
}

func (b *AbstractItem) Decr() (int64, bool) {
	return 0, false
}

func (b *AbstractItem) Xd() (int64, bool) {
	return 0, false
}
