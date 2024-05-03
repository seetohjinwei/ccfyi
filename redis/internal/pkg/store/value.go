package store

type Value struct {
	item Item
}

func NewValue(item Item) *Value {
	ret := &Value{
		item: item,
	}

	return ret
}

func (v *Value) Item() Item {
	return v.item
}
