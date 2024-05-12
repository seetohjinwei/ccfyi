package store

type AbstractItem struct{}

// don't define `Serialise` (all structs should define this)

func (b *AbstractItem) Get() (string, bool) {
	return "", false
}

func (b *AbstractItem) Incr() (int64, bool) {
	return 0, false
}

func (b *AbstractItem) Decr() (int64, bool) {
	return 0, false
}

func (b *AbstractItem) LPush(strs []string) (int64, bool) {
	return 0, false
}

func (b *AbstractItem) RPush(strs []string) (int64, bool) {
	return 0, false
}

func (b *AbstractItem) LRange(start, stop int) ([]string, bool) {
	return []string{}, false
}

func (b *AbstractItem) LLen() (int64, bool) {
	return 0, false
}
