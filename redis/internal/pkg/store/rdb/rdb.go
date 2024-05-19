package rdb

import (
	"bytes"
	"errors"
	"time"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
	"github.com/seetohjinwei/ccfyi/redis/pkg/delay"
)

// *partially* follows the file format listed here: https://rdb.fnordig.de/file_format.html
// liberties taken because this is a Go program, some of these optimisations don't matter

const redis = "REDIS"
const version = "LITE" // intentionally not an integer to not collide with redis version numbers
const magicString = redis + version

type SaveBuffer struct {
	bytes.Buffer
}

func (buf *SaveBuffer) header() {
	buf.WriteString(magicString)
}

func (buf *SaveBuffer) value(k string, v *store.Value) {
	item, ok := v.Item()
	if !ok {
		return
	}

	buf.Write(v.SerialiseExpiry())
	buf.WriteByte(byte(item.ValueType()))
	buf.Write(encoding.EncodeString(k))
	buf.Write(item.Serialise())
}

func (buf *SaveBuffer) checksum() {
	checksum := encoding.GenerateChecksum(buf.Bytes())
	buf.Write(checksum)
}

func (buf *SaveBuffer) eof() {
	buf.WriteString("FF")
}

// Make sure to lock the map!
func (buf *SaveBuffer) Save(values map[string]*store.Value) []byte {
	buf.header()

	for k, v := range values {
		buf.value(k, v)
	}

	buf.eof()
	buf.checksum()

	return buf.Bytes()
}

// zero value is NOT usable.
type LoadBuffer struct {
	b    []byte
	full []byte
	ret  map[string]*store.Value
}

func NewLoadBuffer(b []byte) LoadBuffer {
	return LoadBuffer{
		full: b,
		b:    b,
		ret:  map[string]*store.Value{},
	}
}

func (buf *LoadBuffer) header() error {
	var found bool
	buf.b, found = bytes.CutPrefix(buf.b, []byte(magicString))
	if !found {
		return errors.New("magic string not found")
	}
	return nil
}

func (buf *LoadBuffer) expiry() (*delay.Delay, error) {
	var found bool
	var usec int64
	var err error
	buf.b, found = bytes.CutPrefix(buf.b, []byte("FD"))
	if !found {
		return nil, nil
	}
	usec, buf.b, err = encoding.DecodeInteger(buf.b)
	if err != nil {
		return nil, err
	}
	return delay.NewDelay(time.UnixMicro(usec)), nil
}

func (buf *LoadBuffer) valueType() (encoding.ValueType, error) {
	if len(buf.b) == 0 {
		return 0, errors.New("not enough bytes for valueType()")
	}

	valueType := buf.b[0]

	ret, err := encoding.GetValueType(valueType)
	if err != nil {
		return ret, err
	}

	buf.b = buf.b[1:]
	return ret, err
}

func (buf *LoadBuffer) key() (string, error) {
	var ret string
	var err error

	ret, buf.b, err = encoding.DecodeString(buf.b)

	return ret, err
}

func (buf *LoadBuffer) value(valueType encoding.ValueType) (store.Item, error) {
	var item store.Item
	var err error

	switch valueType {
	case encoding.ValueString:
		item, buf.b, err = store.DeserialiseString(buf.b)
		return item, err
	case encoding.ValueList:
		item, buf.b, err = store.DeserialiseList(buf.b)
		return item, err
	}

	return nil, errors.New("cannot deserialise value because value type is unknown")
}

func (buf *LoadBuffer) item() error {
	expiry, err := buf.expiry()
	if err != nil {
		return err
	}
	valueType, err := buf.valueType()
	if err != nil {
		return err
	}
	key, err := buf.key()
	if err != nil {
		return err
	}
	value, err := buf.value(valueType)
	if err != nil {
		return err
	}

	buf.ret[key] = store.NewValue(value, expiry)

	return nil
}

func (buf *LoadBuffer) values() error {
	for {
		if done, err := buf.eof(); done {
			return err
		}
		if err := buf.item(); err != nil {
			return err
		}
	}
}

func (buf *LoadBuffer) checksum() bool {
	end := len(buf.full) - len(buf.b)
	return encoding.VerifyChecksum(buf.full[:end], buf.b)
}

func (buf *LoadBuffer) eof() (done bool, err error) {
	var found bool
	buf.b, found = bytes.CutPrefix(buf.b, []byte("FF"))
	if found {
		if !buf.checksum() {
			return true, errors.New("checksum did not match")
		}
		return true, nil
	}
	return false, nil
}

func (buf *LoadBuffer) Load() (map[string]*store.Value, error) {
	if err := buf.header(); err != nil {
		return nil, err
	}
	if err := buf.values(); err != nil {
		return nil, err
	}
	return buf.ret, nil
}
