package rdb

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
)

func TestEncodeLength(t *testing.T) {
	Equal(t, V(EncodeLength(0)), V([]byte{0}))
	Equal(t, V(EncodeLength(63)), V([]byte{63}))
	Equal(t, V(EncodeLength(64)), V([]byte{0b01000000, 64}))
	Equal(t, V(EncodeLength(255)), V([]byte{0b01000000, 255}))
	Equal(t, V(EncodeLength(256)), V([]byte{0b01000001, 0}))
	Equal(t, V(EncodeLength(16383)), V([]byte{0b01111111, 255}))
	Equal(t, V(EncodeLength(16384)), V([]byte{0b10000000, 0, 0, 0b1000000, 0}))
	Equal(t, V(EncodeLength(53189571)), V([]byte{0b10000000, 0b11, 0b00101011, 0b10011011, 0b11000011}))
}

func TestDecodeLength(t *testing.T) {
	Equal(t, V(DecodeLength([]byte{0})), V(uint(0), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{63})), V(uint(63), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{0b01000000, 64})), V(uint(64), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{0b01000000, 255})), V(uint(255), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{0b01000001, 0})), V(uint(256), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{0b01111111, 255})), V(uint(16383), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{0b10000000, 0, 0, 0b1000000, 0})), V(uint(16384), []byte{}, nil))
	Equal(t, V(DecodeLength([]byte{0b10000000, 0b11, 0b00101011, 0b10011011, 0b11000011})), V(uint(53189571), []byte{}, nil))
}
